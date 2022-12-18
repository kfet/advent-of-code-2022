package main

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/samber/lo"
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

func (m *mesh) buildDistanceLimitMatrix() map[*valve]map[*valve]int {
	vvm := map[*valve]map[*valve]int{}
	for _, v := range *m {
		vvm[v] = m.buildDistanceLimitRow(v)
	}
	return vvm
}

func unvisitedMaximumLimit(actors []actor, distances map[*valve]map[*valve]int, unvisited map[*valve]struct{}) int {
	var maxLim int
	for cv := range unvisited {
		var maxRate int
		for _, ac := range actors {
			cvTimeLeft := ac.timeLeft - distances[ac.v][cv]
			if cvTimeLeft < 0 {
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

func maxFlow(actors []actor, distances map[*valve]map[*valve]int, unopen map[*valve]struct{}) int {
	// enum order of opening, and calculate distance
	var maxFlowRate int
	for vc := range unopen {
		var maxVcTimeLeft int
		maxTimeActorIndexes := []int{}
		// pick the actor which can get to the valve the soonest
		for i, ac := range actors {
			vcTimeLeft := ac.timeLeft - distances[ac.v][vc]
			if vcTimeLeft > maxVcTimeLeft {
				maxVcTimeLeft = vcTimeLeft
				maxTimeActorIndexes = []int{i}
			} else if vcTimeLeft == maxVcTimeLeft {
				maxTimeActorIndexes = append(maxTimeActorIndexes, i)
			}
		}

		for _, i := range maxTimeActorIndexes {
			// .. for each actor which can get at the earlies time to the valve
			actorsCopy := lo.Map(actors, func(item actor, index int) actor { return item })

			actorsCopy[i].timeLeft -= distances[actors[i].v][vc] + 1
			actorsCopy[i].v = vc

			nextUnopen := input.CopyMap(unopen, func(*valve, struct{}) bool { return true })
			delete(nextUnopen, vc)

			childFlowRate := actorsCopy[i].timeLeft * vc.rate
			maxLim := unvisitedMaximumLimit(actorsCopy, distances, nextUnopen)
			if childFlowRate+maxLim <= maxFlowRate {
				continue
			}

			childFlowRate += maxFlow(actorsCopy, distances, nextUnopen)
			if childFlowRate > maxFlowRate {
				maxFlowRate = childFlowRate
			}
		}
	}

	return maxFlowRate
}

func (m *mesh) maxFlow(actorNames []string, timeLeft int, distances map[*valve]map[*valve]int) int {
	allChildren := lo.MapEntries(*m, func(key string, value *valve) (*valve, struct{}) {
		return value, struct{}{}
	})
	allChildren = input.CopyMap(allChildren, func(v *valve, s struct{}) bool {
		return v.rate > 0
	})

	actors := lo.Map(actorNames, func(item string, index int) actor {
		return actor{
			name:     item,
			v:        (*m)["AA"],
			timeLeft: timeLeft,
		}
	})

	maxFlowRate := maxFlow(actors, distances, allChildren)

	return maxFlowRate
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
}
