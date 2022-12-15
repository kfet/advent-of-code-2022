package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"kfet.org/aoc_common/assert"
	"kfet.org/aoc_common/calc"
	"kfet.org/aoc_common/input"
)

type sensor struct {
	x, y int
	dist int
}

type beacon struct {
	x, y int
}

func NewSensor(x, y, bx, by int) *sensor {
	return &sensor{
		x:    x,
		y:    y,
		dist: calc.TaxiCab(x, y, bx, by),
	}
}

func (s *sensor) rowCoverage(row int) (*rowCoverage, bool) {
	dx := s.dist - calc.Abs(row-s.y)
	if dx < 0 {
		return nil, false
	}
	return &rowCoverage{x1: s.x - dx, x2: s.x + dx}, true
}

func (s *sensor) pointCovered(x, y int) bool {
	return calc.TaxiCab(s.x, s.y, x, y) <= s.dist
}

// Parse string in the format "x=NN, y=NN"
func parseCoords(coords string) (int, int, error) {
	coordsTokens := strings.Split(coords, ", ")
	if len(coordsTokens) != 2 {
		return 0, 0, errors.New("wrong coords format " + coords)
	}

	x, err := strconv.Atoi(coordsTokens[0][len("x="):])
	if err != nil {
		return 0, 0, err
	}
	y, err := strconv.Atoi(coordsTokens[1][len("y="):])
	if err != nil {
		return 0, 0, err
	}

	return x, y, nil
}

// Parse a string in the format "Sensor at <coords>: closest beacon is at <coords>"
// see 'parseCoords' for <coords> format
func readSensor(line string) (*sensor, *beacon, error) {
	tokens := strings.Split(line, ": ")
	if len(tokens) != 2 {
		return nil, nil, errors.New("wrong sensor format " + line)
	}

	x, y, err := parseCoords(tokens[0][len("Sensor at "):])
	if err != nil {
		return nil, nil, err
	}

	bx, by, err := parseCoords(tokens[1][len("closest beacon is at "):])
	if err != nil {
		return nil, nil, err
	}

	return NewSensor(x, y, bx, by), &beacon{bx, by}, nil
}

type rowCoverage struct {
	x1, x2 int
}

func (c *rowCoverage) len() int {
	return c.x2 - c.x1 + 1
}

func (c *rowCoverage) merge(other *rowCoverage) (*rowCoverage, bool) {
	if c.x2 < other.x1-1 || c.x1 > other.x2+1 {
		return nil, false
	}

	res := &rowCoverage{
		x1: calc.Min(c.x1, other.x1),
		x2: calc.Max(c.x2, other.x2),
	}
	return res, true
}

func processFile(name string, interestingRow int, searchSize int, partOne bool) (int, error) {

	// coverage in interesting row
	m := map[*rowCoverage]struct{}{}
	// beacons detected on the interesting row
	mb := map[int]struct{}{}

	ms := map[*sensor]struct{}{}

	err := input.ReadFileLines(name, func(line string) error {
		s, b, err := readSensor(line)
		if err != nil {
			return err
		}

		if b.y == interestingRow {
			mb[b.x] = struct{}{}
		}

		ms[s] = struct{}{}

		rc, cover := s.rowCoverage(interestingRow)
		if !cover {
			return nil
		}

		// merge coverage. That way we can count only once each point in the row
		for c := range m {
			if nc, ok := c.merge(rc); ok {
				*rc = *nc
				delete(m, c)
			}
		}
		m[rc] = struct{}{}

		return nil
	})
	if err != nil {
		return 0, err
	}

	if partOne {
		var sum int
		for c := range m {
			sum += c.len()
		}

		return sum - len(mb), nil
	}

	// part two
	isCovered := func(x, y int) (*sensor, bool) {
		for s := range ms {
			if s.pointCovered(x, y) {
				return s, true
			}
		}
		return nil, false
	}

	for y := 0; y < searchSize; y++ {
		if y%1_000_000 == 0 {
			fmt.Printf("Row %d ..\n", y)
		}
		for x := 0; x < searchSize; x++ {
			if s, isCovered := isCovered(x, y); isCovered {
				coverage, _ := s.rowCoverage(y)
				// skip to the edge of this sensor coverage on the current row
				// ^ this is the important optimizatio to solve the puzzle ^
				x = coverage.x2
				continue
			}
			return x*4_000_000 + y, nil
		}
	}

	return 0, errors.New("beacon not found")
}

func main() {
	res, err := processFile("data/part_one.txt", 10, 20, true)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
	fmt.Println("=================")
	assert.Equals(26, res, "")

	res, err = processFile("data/input.txt", 2_000_000, 4_000_000, true)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
	fmt.Println("=================")
	assert.Equals(4737567, res, "")

	res, err = processFile("data/part_one.txt", 10, 20, false)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
	fmt.Println("=================")
	assert.Equals(56000011, res, "")

	res, err = processFile("data/input.txt", 2_000_000, 4_000_000, false)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
	fmt.Println("=================")
	assert.Equals(13267474686239, res, "")
}
