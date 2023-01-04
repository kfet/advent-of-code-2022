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
	id        int
	robotCost []robotCost // indexed by material
	robotMax  []int       // max number of each robot type, indexed by material
}

type robotCost []int // indexed by material

type worldState struct {
	blueprint *blueprint
	robots    []int
	goods     []int
}

func NewState(blueprintLine string) *worldState {
	s := &worldState{
		blueprint: parseBlueprint(blueprintLine),
		robots:    make([]int, materialsCount),
		goods:     make([]int, materialsCount),
	}
	s.robots[ore] = 1
	return s
}

func (ws *worldState) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintln(ws.blueprint))
	sb.WriteString(fmt.Sprintln("r: ", ws.robots))
	sb.WriteString(fmt.Sprintln("m: ", ws.goods))
	return sb.String()
}

func (ws *worldState) copyState() *worldState {
	return &worldState{
		blueprint: ws.blueprint,
		robots:    input.CopySlice(ws.robots),
		goods:     input.CopySlice(ws.goods),
	}
}

func (ws *worldState) robotBuildTime(rCost robotCost) (int, bool) {
	waitTime := 0
	for mat, cost := range rCost {
		if cost <= ws.goods[mat] {
			// no wait time for this building block
			continue
		}
		// cost > ws.goods[m], i.e. need to wait for material to be produced...
		if ws.robots[mat] == 0 {
			//.. but no robots  esist to build the base materials - return false
			return 0, false
		}

		lacking := cost - ws.goods[mat]
		wait := lacking / ws.robots[mat]
		if lacking%ws.robots[mat] > 0 {
			wait++
		}

		waitTime = calc.Max(waitTime, wait)
	}

	// add 1 minute to build the bot
	return waitTime + 1, true
}

func (ws *worldState) buildBot(mat int, rc robotCost) {
	for m, c := range rc {
		ws.goods[m] -= c
	}
	ws.robots[mat]++
}

func (ws *worldState) mine(time int) {
	for m, count := range ws.robots {
		ws.goods[m] += time * count
	}
}

type nextState struct {
	ws       *worldState
	timeLeft int
}

func (ns *nextState) String() string {
	return fmt.Sprint("tl:", ns.timeLeft, " ws:", ns.ws)
}

func (ws *worldState) nextStates(timeLeft int, target material) []*nextState {
	// choices for next steps: which bot to build next
	res := []*nextState{}

	for mat, rCost := range ws.blueprint.robotCost {
		if material(mat) != target &&
			ws.robots[mat] >= ws.blueprint.robotMax[mat] {
			// already saturated this type of robot
			// skip state
			continue
		}

		waitAndBuildTime, ok := ws.robotBuildTime(rCost)
		if !ok {
			// can't build bot, lacks base material bots
			// skip state
			continue
		}

		if waitAndBuildTime > timeLeft {
			// not enough time to build this type of bot
			// just mine until end of time
			endWs := ws.copyState()
			endWs.mine(timeLeft)
			res = append(res, &nextState{
				ws:       endWs,
				timeLeft: 0,
			})
			continue
		}

		nws := ws.copyState()
		nws.mine(waitAndBuildTime)
		nws.buildBot(mat, rCost)
		res = append(res, &nextState{
			ws:       nws,
			timeLeft: timeLeft - waitAndBuildTime,
		})
	}

	return res
}

func (ws *worldState) upperLimitGoods(timeLeft int, mat material) int {
	return ws.goods[mat] +
		timeLeft*ws.robots[mat] +
		timeLeft*timeLeft/2
}

func (ws *worldState) maxGoods(timeLeft int, target material, minGoods int) (int, []*nextState) {
	if timeLeft == 0 {
		// at end of search, return goods produced
		return ws.goods[target], []*nextState{{ws: ws, timeLeft: timeLeft}}
	}

	if ws.upperLimitGoods(timeLeft, target) < minGoods {
		// theoretical maximum lower than required minimum
		return 0, nil
	}

	max := minGoods
	var maxWs []*nextState
	for _, ns := range ws.nextStates(timeLeft, target) {
		nextMax, mws := ns.ws.maxGoods(ns.timeLeft, target, max)
		if nextMax > max {
			max = nextMax
			maxWs = mws
		}
	}

	return max, append([]*nextState{{ws: ws, timeLeft: timeLeft}}, maxWs...)
}

// "Blueprint 1: Each ore robot costs 4 ore. Each clay robot costs 2 ore. Each obsidian robot costs 3 ore and 14 clay. Each geode robot costs 2 ore and 7 obsidian."
func parseBlueprint(line string) *blueprint {
	tokens := strings.Split(line, " ")
	bp := &blueprint{
		id: input.MustAtoi(tokens[1][0 : len(tokens[1])-1]),
		robotCost: []robotCost{
			make([]int, materialsCount), // ore
			make([]int, materialsCount), // clay
			make([]int, materialsCount), // obsidian
			make([]int, materialsCount), // geode
		},
		robotMax: make([]int, materialsCount),
	}

	bp.robotCost[ore][ore] = input.MustAtoi(tokens[6])

	bp.robotCost[clay][ore] = input.MustAtoi(tokens[12])

	bp.robotCost[obsidian][ore] = input.MustAtoi(tokens[18])
	bp.robotCost[obsidian][clay] = input.MustAtoi(tokens[21])

	bp.robotCost[geode][ore] = input.MustAtoi(tokens[27])
	bp.robotCost[geode][obsidian] = input.MustAtoi(tokens[30])

	for m := 0; m < materialsCount; m++ {
		bp.robotMax[m] = bp.maxRobotsForMaterial(material(m))
	}

	return bp
}

func (b *blueprint) maxRobotsForMaterial(m material) int {
	max := 0
	for _, rc := range b.robotCost {
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
		s := NewState(line)
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

			max, _ := s.maxGoods(timeToRun, mat, 0)
			fmt.Println("max: ", max)

			totalQ += s.blueprint.id * max
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

		max, _ := s.maxGoods(timeToRun, mat, 0)
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
