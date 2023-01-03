package main

import (
	"fmt"
	"strings"
	"time"

	"kfet.org/aoc_common/assert"
	"kfet.org/aoc_common/calc"
	"kfet.org/aoc_common/input"
)

type material uint8

const (
	ore = iota
	clay
	obsidian
	geode
	materialsCount
)

type blueprint struct {
	id    int
	rCost []robotCost
	rMax  []int // max number of each robot type
}

type robotCost []int

type worldState struct {
	b      *blueprint
	robots []int
	goods  []int
}

func NewSim(blueprintLine string) *worldState {
	s := &worldState{
		b:      parseBlueprint(blueprintLine),
		goods:  make([]int, materialsCount),
		robots: make([]int, materialsCount),
	}
	s.robots[ore] = 1
	return s
}

func (ws *worldState) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintln(ws.b))
	sb.WriteString(fmt.Sprintln("r: ", ws.robots))
	sb.WriteString(fmt.Sprint("m: ", ws.goods))
	return sb.String()
}

func (ws *worldState) copySim() *worldState {
	return &worldState{
		b:      ws.b,
		robots: input.CopySlice(ws.robots),
		goods:  input.CopySlice(ws.goods),
	}
}

func (ws *worldState) canBuildBot(rc robotCost) bool {
	for m, c := range rc {
		if ws.goods[m] < c {
			return false
		}
	}
	return true
}

func (ws *worldState) buildBot(rc robotCost) {
	for m, c := range rc {
		ws.goods[m] -= c
	}
}

func (ws *worldState) nextStates(target material) []*worldState {
	// choices for next step: just mine, or build any of the robots
	res := []*worldState{}

	// just mining - each robot we have would produce 1 unit of the material it mines
	justMining := ws.copySim()
	for m, c := range ws.robots {
		justMining.goods[m] += c
	}

	// enum in reverse, giving priority to creating geode-cracking bots
	for m := materialsCount - 1; m >= 0; m-- {
		if material(m) != target &&
			ws.robots[m] >= ws.b.rMax[m] {
			// already saturated this type of robot, skip state
			continue
		}

		rc := ws.b.rCost[m]
		if !ws.canBuildBot(rc) {
			continue
		}

		// clone state from the mining one, and build a bot in it
		ns := justMining.copySim()
		ns.buildBot(rc)
		ns.robots[m]++

		res = append(res, ns)
	}

	// append the waiting state at the very end as an optimization
	res = append(res, justMining)

	return res
}

func (ws *worldState) upperLimitGoods(t int, m material) int {
	rc := ws.b.rCost[m]
	timeToNextMRobo := 1
	for m, c := range rc {
		mcMatRobots := ws.robots[m]
		mcMatGoods := ws.goods[m]
		mcCostMissing := c - mcMatGoods
		t := 1
		for mcCostMissing > 0 {
			mcCostMissing -= mcMatRobots
			mcMatRobots++ // assume best case where we can build those robots each minute
			t++
		}
		timeToNextMRobo = calc.Max(timeToNextMRobo, t)
	}

	return ws.goods[m] +
		t*ws.robots[m] +
		(t-timeToNextMRobo)*(t-timeToNextMRobo)/2
}

func (ws *worldState) maxGoods(time int, target material, minGoods int) int {
	if time == 0 {
		// at end of search, return goods produced
		return ws.goods[target]
	}

	if ws.upperLimitGoods(time, target) < minGoods {
		// theoretical maximum lower than required minimum
		return 0
	}

	maxGoods := minGoods
	for _, ns := range ws.nextStates(target) {
		mg := ns.maxGoods(time-1, target, maxGoods)
		if mg > maxGoods {
			maxGoods = mg
		}
	}

	return maxGoods
}

// "Blueprint 1: Each ore robot costs 4 ore. Each clay robot costs 2 ore. Each obsidian robot costs 3 ore and 14 clay. Each geode robot costs 2 ore and 7 obsidian."
func parseBlueprint(line string) *blueprint {
	tokens := strings.Split(line, " ")
	bp := &blueprint{
		id: input.MustAtoi(tokens[1][0 : len(tokens[1])-1]),
		rCost: []robotCost{
			make([]int, materialsCount), // ore
			make([]int, materialsCount), // clay
			make([]int, materialsCount), // obsidian
			make([]int, materialsCount), // geode
		},
		rMax: make([]int, materialsCount),
	}

	bp.rCost[ore][ore] = input.MustAtoi(tokens[6])

	bp.rCost[clay][ore] = input.MustAtoi(tokens[12])

	bp.rCost[obsidian][ore] = input.MustAtoi(tokens[18])
	bp.rCost[obsidian][clay] = input.MustAtoi(tokens[21])

	bp.rCost[geode][ore] = input.MustAtoi(tokens[27])
	bp.rCost[geode][obsidian] = input.MustAtoi(tokens[30])

	for i := 0; i < materialsCount; i++ {
		max := bp.maxRobotsForMaterial(material(i))
		bp.rMax[i] = max
	}

	return bp
}

func (b *blueprint) maxRobotsForMaterial(m material) int {
	max := 0
	for _, rc := range b.rCost {
		mrc := rc[m]
		if mrc > max {
			max = mrc
		}
	}
	return max
}

func processFile(fileName string, mat material, timeToRun int, partOne bool) (int, error) {
	// read all blueprints into world states
	states := []*worldState{}
	err := input.ReadFileLines(fileName, func(line string) error {
		s := NewSim(line)
		states = append(states, s)
		return nil
	})
	if err != nil {
		return 0, err
	}

	if partOne {
		// run each state
		var totalQ int
		for _, s := range states {
			start := time.Now()
			fmt.Println("Starting at ", start)

			fmt.Println(s.String())

			max := s.maxGoods(timeToRun, mat, 0)
			fmt.Println("max: ", max)

			totalQ += s.b.id * max
			fmt.Println("elapsed time: ", time.Since(start))
			fmt.Println("---------------------------")
		}

		return totalQ, nil
	}

	// part two
	res := 1
	for i := 0; i < 3; i++ {
		s := states[i]

		start := time.Now()
		fmt.Println("Starting at ", start)

		fmt.Println(s.String())

		max := s.maxGoods(timeToRun, mat, 0)
		fmt.Println("max: ", max)

		res *= max

		fmt.Println("elapsed time: ", time.Since(start))
		fmt.Println("---------------------------")
	}

	return res, nil
}

func main() {
	var res int
	var err error
	res, err = processFile("data/part_one.txt", geode, 24, true)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
	fmt.Println("=================")
	assert.Equals(33, res, "")

	// NOTE: this takes a few minutes to finish
	res, err = processFile("data/input.txt", geode, 24, true)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
	fmt.Println("=================")
	assert.Equals(1346, res, "")

	res, err = processFile("data/input.txt", geode, 32, false)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
	fmt.Println("=================")
	assert.Equals(7644, res, "")
}
