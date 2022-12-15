package main

import (
	"errors"
	"fmt"
	"math"
	"strings"

	"kfet.org/aoc_common/assert"
	"kfet.org/aoc_common/input"
)

type dot uint8

const (
	air dot = iota
	rock
	sand
)

type point struct {
	x, y int
}

type caveMap struct {
	m          [][]dot
	off        offset
	maxY       int
	floor      bool
	sandOrigin point
}

func (cm caveMap) String() string {
	var sb strings.Builder
	for y := range cm.m {
		for x := range cm.m[y] {
			switch cm.m[y][x] {
			case air:
				sb.WriteRune('.')
			case rock:
				sb.WriteRune('#')
			case sand:
				sb.WriteRune('o')
			}
		}
		sb.WriteRune('\n')
	}
	return sb.String()
}

func NewCaveMap(name string, s size, off offset, floor bool) *caveMap {
	cm := new(caveMap)

	cm.floor = floor
	if floor {
		// add room for a floor
		s.h += 2

		// expand width to cover twice the height, with room to spare
		// that way we can collect all sand on the floor
		dw := 2*(s.h+10) - s.w
		if dw < 0 {
			dw = 0
		}
		s.w += dw

		// fix X offset accordingly
		off.x += dw / 2
	}

	cm.off = off
	cm.sandOrigin = cm.off.offsetPoint(point{x: 500, y: 0})

	// allocate the area matrix
	cm.m = make([][]dot, s.h)
	for i := range cm.m {
		cm.m[i] = make([]dot, s.w)
	}

	// load from file
	cm.readMap(name)

	return cm
}

func (cm *caveMap) drawLine(from, to point, d dot) error {
	from = cm.off.offsetPoint(from)
	to = cm.off.offsetPoint(to)

	if from.y > cm.maxY {
		cm.maxY = from.y
	}
	if to.y > cm.maxY {
		cm.maxY = from.y
	}

	switch {
	case from.x == to.x:
		// vertical line
		step := 1
		if from.y > to.y {
			step = -1
		}
		for y := from.y; y != to.y; y += step {
			cm.m[y][from.x] = d
		}
		cm.m[to.y][to.x] = d
	case from.y == to.y:
		// horizontal line
		step := 1
		if from.x > to.x {
			step = -1
		}
		for x := from.x; x != to.x; x += step {
			cm.m[from.y][x] = d
		}
		cm.m[to.y][to.x] = d
	default:
		return errors.New("diagonal lines not supported")
	}
	return nil
}

func parsePoints(line string, pointHandler func(point) error) error {
	strPoints := strings.Split(line, " -> ")

	for _, strPoint := range strPoints {
		tokens := strings.Split(strPoint, ",")
		ints := input.MustAtoInts(tokens)
		err := pointHandler(point{x: ints[0], y: ints[1]})
		if err != nil {
			return err
		}
	}
	return nil
}

func (cm *caveMap) readMap(name string) error {
	err := input.ReadFileLines(name, func(line string) error {

		first := true
		var prevPoint point
		parsePoints(line, func(p point) error {
			if first {
				// skip drawing a line on the first point
				first = false
			} else {
				err := cm.drawLine(prevPoint, p, rock)
				if err != nil {
					return err
				}
			}
			prevPoint = p
			return nil
		})

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (cm caveMap) fall(p *point) bool {
	if p.y >= len(cm.m)-2 && cm.floor {
		// stuck to permanent floor
		return false
	}

	if p.y >= len(cm.m)-1 ||
		cm.m[p.y+1][p.x] == air {
		// free fall
		(*p).y++
		return true
	}

	if p.x <= 0 ||
		cm.m[p.y+1][p.x-1] == air {
		// fall left
		(*p).x--
		(*p).y++
		return true
	}

	if p.x >= len(cm.m[p.y+1])-1 ||
		cm.m[p.y+1][p.x+1] == air {
		// fall right
		(*p).x++
		(*p).y++
		return true
	}

	// stuck
	return false
}

func (cm *caveMap) runSand() int {
	var sandUnits int
	for {
		sandGrain := cm.sandOrigin
		falling := cm.fall(&sandGrain)
		for falling {
			// falling...
			if sandGrain.x < 0 || sandGrain.x > len(cm.m[0]) || sandGrain.y > len(cm.m) {
				// sand fell out of the map
				if cm.floor {
					// floor is enabled, should not be reachable
					panic("sand fell out of the map " + fmt.Sprint(sandGrain))
				}
				// finish
				return sandUnits
			}
			falling = cm.fall(&sandGrain)
		}

		// sand grain is stuck
		cm.m[sandGrain.y][sandGrain.x] = sand
		sandUnits++

		if sandGrain == cm.sandOrigin {
			// stuck at the sand entrance, stop the sand
			return sandUnits
		}
	}
}

type size struct {
	w, h int
}

type offset struct {
	x, y int
}

func (off offset) offsetPoint(p point) point {
	return point{
		x: p.x + off.x,
		y: p.y + off.y,
	}
}

func preProcessFile(name string) (point, point, error) {
	minX := math.MaxInt
	var maxX, maxY int

	err := input.ReadFileLines(name, func(line string) error {
		parsePoints(line, func(p point) error {
			if p.x < minX {
				minX = p.x
			}
			if p.x > maxX {
				maxX = p.x
			}
			if p.y > maxY {
				maxY = p.y
			}
			return nil
		})
		return nil
	})
	if err != nil {
		return point{}, point{}, err
	}

	return point{minX, 0}, point{maxX, maxY}, nil
}

func processFile(name string, floor bool) (int, error) {

	// extract map area first
	minP, maxP, err := preProcessFile(name)
	if err != nil {
		return 0, err
	}
	s := size{
		w: maxP.x - minP.x + 1,
		h: maxP.y + 1,
	}
	off := offset{x: -minP.x}

	// load the map and run the sand
	cm := NewCaveMap(name, s, off, floor)
	res := cm.runSand()

	return res, nil
}

func main() {
	res, err := processFile("data/part_one.txt", false)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
	fmt.Println("=================")
	assert.Equals(24, res, "")

	res, err = processFile("data/input.txt", false)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
	fmt.Println("=================")
	assert.Equals(745, res, "")

	res, err = processFile("data/part_one.txt", true)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
	fmt.Println("=================")
	assert.Equals(93, res, "")

	res, err = processFile("data/input.txt", true)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
	fmt.Println("=================")
	assert.Equals(27551, res, "")
}
