package main

import (
	"fmt"
	"strings"

	"kfet.org/aoc_common/input"
)

const timeToRun int = 24 // minutes

type material uint8

const (
	ore = iota
	clay
	obsidian
	geode
)

type blueprint struct {
	id int
	rc []*robotCost
}

func (b *blueprint) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintln(b.id))
	for _, rc := range b.rc {
		sb.WriteString("; " + fmt.Sprint(rc.m) + ": [")
		for _, mc := range rc.cost {
			sb.WriteString(fmt.Sprint(mc.m) + ":" + fmt.Sprint(mc.c) + ", ")
		}
		sb.WriteString("]")
	}
	return sb.String()
}

type robotCost struct {
	m    material
	cost []*materialCost
}

type materialCost struct {
	m material
	c int
}

type robot struct {
	m material
}

type simulation struct {
	b      *blueprint
	robots []*robot
	goods  map[material]int
}

func NewSim(blueprintLine string) *simulation {
	b := parseBlueprint(blueprintLine)
	goods := map[material]int{}
	for _, rc := range b.rc {
		goods[rc.m] = 0
	}
	return &simulation{
		b:      b,
		robots: []*robot{{m: ore}},
		goods:  goods,
	}
}

func (sim *simulation) String() string {
	var sb strings.Builder
	sb.WriteString(sim.b.String())
	sb.WriteString(fmt.Sprintln(sim.robots))
	sb.WriteString(fmt.Sprint(sim.goods))
	return sb.String()
}

func (sim *simulation) copySim() *simulation {
	return &simulation{
		b:      sim.b,
		robots: input.CopySlice(sim.robots),
		goods:  input.CopyMap(sim.goods, input.NoFilter[material, int]),
	}
}

func (sim *simulation) buildBot(rc *robotCost) (*robot, bool) {
	for _, mc := range rc.cost {
		sim.goods[mc.m] -= mc.c
		if sim.goods[mc.m] < 0 {
			return nil, false
		}
	}
	return &robot{rc.m}, true
}

func (sim *simulation) mineGoods() {
	for _, r := range sim.robots {
		sim.goods[r.m]++
	}
}

func (sim *simulation) maxGoods(timeLeft int, m material) int {
	if timeLeft == 0 {
		return 0
	}

	// run a one minute-step for no bot built
	var maxGoods int

	// run a one minute-step variant for each type of bot
	for _, rc := range sim.b.rc {
		ssim := sim.copySim()
		r, ok := ssim.buildBot(rc)
		if !ok {
			// not enough goods to build this bot, skip
			continue
		}
		// mine materials
		ssim.mineGoods()
		ssim.robots = append(ssim.robots, r)
		subMaxGoods := ssim.maxGoods(timeLeft-1, m)
		if subMaxGoods > maxGoods {
			maxGoods = subMaxGoods
		}
	}

	if maxGoods == 0 {
		// no bot can be built, try just mining goods
		ssim := sim.copySim()
		ssim.mineGoods()
		maxGoods = ssim.maxGoods(timeLeft-1, m)
	}

	return sim.goods[m] + maxGoods
}

// "Blueprint 1: Each ore robot costs 4 ore. Each clay robot costs 2 ore. Each obsidian robot costs 3 ore and 14 clay. Each geode robot costs 2 ore and 7 obsidian."
func parseBlueprint(line string) *blueprint {
	tokens := strings.Split(line, " ")
	return &blueprint{
		id: input.MustAtoi(tokens[1][0 : len(tokens[1])-1]),
		rc: []*robotCost{
			{m: ore, cost: []*materialCost{
				{
					m: ore,
					c: input.MustAtoi(tokens[6]),
				},
			}},
			{m: clay, cost: []*materialCost{
				{
					m: ore,
					c: input.MustAtoi(tokens[12]),
				},
			}},
			{m: obsidian, cost: []*materialCost{
				{
					m: ore,
					c: input.MustAtoi(tokens[18]),
				},
				{
					m: clay,
					c: input.MustAtoi(tokens[21]),
				},
			}},
			{m: geode, cost: []*materialCost{
				{
					m: ore,
					c: input.MustAtoi(tokens[27]),
				},
				{
					m: obsidian,
					c: input.MustAtoi(tokens[30]),
				},
			}},
		},
	}
}

func processFile(fileName string) (int, error) {

	sims := []*simulation{}

	err := input.ReadFileLines(fileName, func(line string) error {
		sims = append(sims, NewSim(line))
		return nil
	})
	if err != nil {
		return 0, err
	}

	var totalQ int
	for _, s := range sims {
		fmt.Println(s.String())
		maxGoods := s.maxGoods(timeToRun, geode)
		fmt.Println(maxGoods)
		totalQ += s.b.id * maxGoods
	}

	return totalQ, nil
}

func main() {
	res, err := processFile("data/part_one.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
	fmt.Println("=================")

	// res, err = processFile("data/input.txt")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println(res)
	// fmt.Println("=================")
}
