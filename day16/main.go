package main

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/samber/lo"
	"kfet.org/aoc_common/assert"
	"kfet.org/aoc_common/calc"
	"kfet.org/aoc_common/input"
)

type valve struct {
	name    string
	rate    int
	state   bool
	tunnels map[*valve]struct{}
}

func NewValve(name string, rate int) *valve {
	return &valve{
		name:    name,
		rate:    rate,
		tunnels: make(map[*valve]struct{}),
	}
}

func (v *valve) String() string {
	var sb strings.Builder
	ts := lo.Map(lo.Keys(v.tunnels), func(item *valve, index int) string {
		return item.name
	})
	sb.WriteString(fmt.Sprintf("name: %s, rate: %d, tunnels: %s\n", v.name, v.rate, strings.Join(ts, ", ")))
	return sb.String()
}

type mesh map[string]*valve

func NewMesh(fileName string) (mesh, error) {
	m := make(mesh)
	err := input.ReadFileLines(fileName, func(line string) error {
		_, err := m.readValve(line)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (m *mesh) String() string {
	var sb strings.Builder
	for _, v := range *m {
		sb.WriteString(v.String())
	}
	return sb.String()
}

func (m *mesh) readValve(line string) (*valve, error) {
	// "Valve AA has flow rate=0; tunnels lead to valves DD, II, BB"
	// "Valve AA has flow rate=0; tunnel leads to valve DD"
	tokens := strings.Split(line, " ")
	if len(tokens) < 10 {
		return nil, errors.New("invalid valve format " + line)
	}

	name := tokens[1]
	rate := input.MustAtoi(tokens[4][5 : len(tokens[4])-1])

	var (
		v     *valve
		found bool
	)
	if v, found = (*m)[name]; !found {
		v = NewValve(name, rate)
		(*m)[name] = v
	} else {
		v.rate = rate
	}

	lo.ForEach(tokens[9:], func(item string, _ int) {
		if item[len(item)-1] == ',' {
			// string last ","
			item = item[0 : len(item)-1]
		}
		var tv *valve
		if tv, found = (*m)[item]; !found {
			tv = NewValve(item, -1)
			(*m)[item] = tv
		}
		v.tunnels[tv] = struct{}{}
	})

	return v, nil
}

func (m *mesh) buildDistanceLimitRow(v *valve) map[*valve]int {
	vm := map[*valve]int{}

	level := 0
	nextWave := []*valve{}
	wave := []*valve{v}

	for len(wave) > 0 {
		for _, vc := range wave {
			if _, visited := vm[vc]; visited {
				// visited, skip
				continue
			}
			vm[vc] = level
			nextWave = append(nextWave, lo.Keys(vc.tunnels)...)
		}
		level++
		wave = nextWave
		nextWave = []*valve{}
	}

	return vm
}

type valveToValveMatrix map[*valve]map[*valve]int

func (m *mesh) buildDistanceLimitMatrix() valveToValveMatrix {
	vvm := map[*valve]map[*valve]int{}
	for _, v := range *m {
		vvm[v] = m.buildDistanceLimitRow(v)
	}
	return vvm
}

func unvisitedUpperRateLimit(actors []actor, distMx valveToValveMatrix, unvisited map[*valve]struct{}) int {
	var maxLim int
	for cv := range unvisited {
		// find max rate for this node
		var maxRate int
		for _, ac := range actors {
			cvTimeLeft := ac.timeLeft - ac.timeToOpenValve(cv, distMx)
			if cvTimeLeft <= 0 {
				continue
			}
			acRate := cv.rate * cvTimeLeft
			if acRate > maxRate {
				maxRate = acRate
			}
		}
		maxLim += maxRate
	}
	return maxLim
}

type actor struct {
	v        *valve
	timeLeft int
	name     string
}

func (a *actor) timeToOpenValve(v *valve, distMx valveToValveMatrix) int {
	return distMx[a.v][v] + 1
}

func maxFlow(actors []actor, distMx valveToValveMatrix, unopen map[*valve]struct{}, requiredFlowRate int) int {
	// enum order of opening, and calculate distance
	var maxFlowRate int
	for nextValve := range unopen {
		var maxNvTimeLeft int
		maxTimeActorIndexes := []int{}
		// pick the actor which can get to the valve the soonest
		for i, ac := range actors {
			nvTimeLeft := ac.timeLeft - ac.timeToOpenValve(nextValve, distMx)
			if nvTimeLeft > maxNvTimeLeft {
				maxNvTimeLeft = nvTimeLeft
				maxTimeActorIndexes = []int{i}
			} else if nvTimeLeft == maxNvTimeLeft {
				maxTimeActorIndexes = append(maxTimeActorIndexes, i)
			}
		}

		if maxNvTimeLeft <= 0 {
			continue
		}

		for _, i := range maxTimeActorIndexes {
			// for each actor which can get at the earliest time to the valve
			actorsCopy := lo.Map(actors, func(item actor, index int) actor { return item })

			// move selected actor to next valve
			actorsCopy[i].timeLeft -= actorsCopy[i].timeToOpenValve(nextValve, distMx)
			actorsCopy[i].v = nextValve

			// consider the next valve open
			nextValveFlowRate := actorsCopy[i].timeLeft * nextValve.rate

			nextUnopen := input.CopyMap(unopen, func(item *valve, s struct{}) bool { return item != nextValve })
			upperLim := unvisitedUpperRateLimit(actorsCopy, distMx, nextUnopen)
			if nextValveFlowRate+upperLim <= maxFlowRate {
				continue
			}
			if nextValveFlowRate+upperLim <= requiredFlowRate {
				continue
			}

			nextValveFlowRate += maxFlow(actorsCopy, distMx, nextUnopen,
				// pass required minimum flow rate to the recursive call
				calc.Max(maxFlowRate-nextValveFlowRate, requiredFlowRate-nextValveFlowRate))
			if nextValveFlowRate > maxFlowRate {
				maxFlowRate = nextValveFlowRate
			}
		}
	}

	return maxFlowRate
}

func (m *mesh) maxFlow(actorNames []string, timeLeft int, distMx valveToValveMatrix) int {
	allChildren := lo.MapEntries(*m, func(key string, value *valve) (*valve, struct{}) {
		return value, struct{}{}
	})
	allChildren = input.CopyMap(allChildren, func(v *valve, s struct{}) bool {
		return v.rate > 0
	})

	actors := lo.Map(actorNames, func(name string, index int) actor {
		return actor{
			name:     name,
			v:        (*m)["AA"],
			timeLeft: timeLeft,
		}
	})

	return maxFlow(actors, distMx, allChildren, 0)
}

func processFile(fileName string, actorNames []string, timeLeft int) (int, error) {
	m, err := NewMesh(fileName)
	if err != nil {
		return 0, err
	}

	matrix := m.buildDistanceLimitMatrix()
	maxFlow := m.maxFlow(actorNames, timeLeft, matrix)

	return maxFlow, nil
}

func main() {
	t := time.Now()
	res, err := processFile("data/part_one.txt", []string{"me"}, 30)
	d := time.Now().Sub(t)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
	fmt.Println("Duration: ", d)
	fmt.Println("=================")
	assert.Equals(1651, res, "")

	t = time.Now()
	res, err = processFile("data/input.txt", []string{"me"}, 30)
	d = time.Now().Sub(t)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
	fmt.Println("Duration: ", d)
	fmt.Println("=================")
	assert.Equals(1659, res, "")

	t = time.Now()
	res, err = processFile("data/part_one.txt", []string{"me", "elephant"}, 26)
	d = time.Now().Sub(t)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
	fmt.Println("Duration: ", d)
	fmt.Println("=================")
	assert.Equals(1707, res, "")

	t = time.Now()
	res, err = processFile("data/input.txt", []string{"me", "elephant"}, 26)
	d = time.Now().Sub(t)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
	fmt.Println("Duration: ", d)
	fmt.Println("=================")
	assert.Equals(2382, res, "")
}
