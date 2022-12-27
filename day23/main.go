package main

import (
	"errors"
	"fmt"
	"math"
	"strings"

	"kfet.org/aoc_common/assert"
	"kfet.org/aoc_common/input"
)

var noMoveRule = rule{
	0, 0,
	[]pos{
		{-1, -1},
		{0, -1},
		{1, -1},
		{-1, 0},
		{1, 0},
		{-1, 1},
		{0, 1},
		{1, 1},
	},
}
var rulesCount int = len(rules)
var rules []rule = []rule{
	{ // N
		0, -1,
		[]pos{
			{-1, -1},
			{0, -1},
			{1, -1},
		},
	},
	{ // S
		0, 1,
		[]pos{
			{-1, 1},
			{0, 1},
			{1, 1},
		},
	},
	{ // W
		-1, 0,
		[]pos{
			{-1, -1},
			{-1, 0},
			{-1, 1},
		},
	},
	{ // E
		1, 0,
		[]pos{
			{1, -1},
			{1, 0},
			{1, 1},
		},
	},
}

type rule struct {
	dx, dy int
	tests  []pos
}

type field struct {
	t int
	m map[int]map[int]struct{}
}

func NewField() *field {
	return &field{
		m: map[int]map[int]struct{}{},
	}
}

func (f *field) String() string {
	var sb strings.Builder
	minp, maxp, _ := f.findBoundaries()
	for y := minp.y; y <= maxp.y; y++ {
		for x := minp.x; x <= maxp.x; x++ {
			if f.isPresent(x, y) {
				sb.WriteRune('#')
			} else {
				sb.WriteRune('.')
			}
		}
		sb.WriteRune('\n')
	}
	return sb.String()
}

func (f *field) findBoundaries() (pos, pos, int) {
	minX, minY := math.MaxInt, math.MaxInt
	maxX, maxY := math.MinInt, math.MinInt
	count := 0
	for y, row := range f.m {
		for x := range row {
			count++
			if x < minX {
				minX = x
			}
			if x > maxX {
				maxX = x
			}
			if y < minY {
				minY = y
			}
			if y > maxY {
				maxY = y
			}
		}
	}

	return pos{minX, minY}, pos{maxX, maxY}, count
}

func (f *field) testRule(x, y int, r rule) bool {
	for _, t := range r.tests {
		if f.isPresent(x+t.x, y+t.y) {
			return false
		}
	}
	return true
}

func (f *field) tick() bool {
	// compile proposed moves
	pm := NewProposedMoves()
	for y, row := range f.m {
		for x := range row {
			// try no-move rule first
			if f.testRule(x, y, noMoveRule) {
				// stay put, no-move-rule matches
				continue
			}

			for i := 0; i < rulesCount; i++ {
				// try each rule
				r := rules[(f.t+i)%rulesCount]
				if !f.testRule(x, y, r) {
					// can't apply rule, try the next one
					continue
				}
				// rule matches
				pm.propose(x, y, r)
				break // .. from rules loop
			}
		}
	}

	// apply proposed moves
	anyMove := false
	for to, from := range pm.toFrom {
		if _, ok := pm.discarded[to]; ok {
			// skip discarded
			continue
		}
		f.move(from, to)
		anyMove = true
	}

	f.t++
	return anyMove
}

func (f *field) isPresent(x, y int) bool {
	if _, ok := f.m[y]; !ok {
		return false
	}
	if _, ok := f.m[y][x]; !ok {
		return false
	}
	return true
}

func (f *field) set(p pos) {
	if _, ok := f.m[p.y]; !ok {
		f.m[p.y] = map[int]struct{}{}
	}
	f.m[p.y][p.x] = struct{}{}
}

func (f *field) unset(p pos) {
	delete(f.m[p.y], p.x)
}

func (f *field) move(from, to pos) {
	f.unset(from)
	f.set(to)
}

type pos struct {
	x, y int
}

type proposedMoves struct {
	toFrom    map[pos]pos
	discarded map[pos]struct{}
}

func NewProposedMoves() *proposedMoves {
	return &proposedMoves{
		toFrom:    map[pos]pos{},
		discarded: map[pos]struct{}{},
	}
}

func (pm *proposedMoves) propose(x, y int, r rule) bool {
	to := pos{x + r.dx, y + r.dy}
	if _, ok := pm.discarded[to]; ok {
		return false
	}

	if _, ok := pm.toFrom[to]; ok {
		pm.discarded[to] = struct{}{}
		return false
	}

	pm.toFrom[to] = pos{x, y}
	return true
}

func processFile(fileName string, partOne bool) (int, error) {

	f := NewField()

	var row int
	err := input.ReadFileLines(fileName, func(line string) error {
		for x, r := range line {
			switch r {
			case '#':
				f.set(pos{x, row})
			case '.':
			default:
				return errors.New(fmt.Sprint("wrong character in map ", r))
			}
		}
		row++
		return nil
	})
	if err != nil {
		return 0, err
	}

	fmt.Println("======================")
	fmt.Println(f)

	var res int
	if partOne {
		for i := 0; i < 10; i++ {
			f.tick()
		}

		minp, maxp, count := f.findBoundaries()
		res = (maxp.x-minp.x+1)*(maxp.y-minp.y+1) - count
	} else {
		for i := 0; ; i++ {
			if !f.tick() {
				res = i + 1
				break
			}
		}
	}

	return res, nil
}

func main() {
	res, err := processFile("data/part_one_small.txt", true)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
	fmt.Println("=================")
	assert.Equals(25, res, "")

	res, err = processFile("data/part_one.txt", true)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
	fmt.Println("=================")
	assert.Equals(110, res, "")

	res, err = processFile("data/input.txt", true)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
	fmt.Println("=================")
	assert.Equals(4082, res, "")

	res, err = processFile("data/part_one.txt", false)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
	fmt.Println("=================")
	assert.Equals(20, res, "")

	res, err = processFile("data/input.txt", false)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
	fmt.Println("=================")
	assert.Equals(1065, res, "")
}
