package main

import (
	"fmt"
	"math"
	"strings"

	"github.com/samber/lo"
	"kfet.org/aoc_common/assert"
	"kfet.org/aoc_common/input"
)

type cube struct {
	x, y, z int
}

func NewCube(ints []int) *cube {
	return &cube{
		x: ints[0],
		y: ints[1],
		z: ints[2],
	}
}

type coordsMap struct {
	minX, minY, minZ int
	maxX, maxY, maxZ int
	m                map[int]map[int]map[int]material
}

func NewCoordsMap() *coordsMap {
	return &coordsMap{
		minX: math.MaxInt,
		minY: math.MaxInt,
		minZ: math.MaxInt,
		maxX: math.MinInt,
		maxY: math.MinInt,
		maxZ: math.MinInt,
		m:    map[int]map[int]map[int]material{},
	}
}

type material uint8

const (
	air = iota
	lava
	pocket_air
)

func (cm *coordsMap) updateMinMax(c *cube) {
	if cm.minX > c.x {
		cm.minX = c.x
	}
	if cm.maxX < c.x {
		cm.maxX = c.x
	}
	if cm.minY > c.y {
		cm.minY = c.y
	}
	if cm.maxY < c.y {
		cm.maxY = c.y
	}
	if cm.minZ > c.z {
		cm.minZ = c.z
	}
	if cm.maxZ < c.z {
		cm.maxZ = c.z
	}
}

func (cm *coordsMap) set(c *cube, m material) {
	cm.updateMinMax(c)
	if _, ok := cm.m[c.x]; !ok {
		cm.m[c.x] = map[int]map[int]material{}
	}
	if _, ok := cm.m[c.x][c.y]; !ok {
		cm.m[c.x][c.y] = map[int]material{}
	}
	cm.m[c.x][c.y][c.z] = m
}

func (cm *coordsMap) get(c *cube) (material, bool) {
	if _, ok := cm.m[c.x]; !ok {
		return 0, false
	}
	if _, ok := cm.m[c.x][c.y]; !ok {
		return 0, false
	}
	if m, ok := cm.m[c.x][c.y][c.z]; !ok {
		return 0, false
	} else {
		return m, true
	}
}

func (c *cube) neighbours() []*cube {
	res := []*cube{}
	res = append(res, &cube{c.x - 1, c.y, c.z})
	res = append(res, &cube{c.x + 1, c.y, c.z})
	res = append(res, &cube{c.x, c.y - 1, c.z})
	res = append(res, &cube{c.x, c.y + 1, c.z})
	res = append(res, &cube{c.x, c.y, c.z - 1})
	res = append(res, &cube{c.x, c.y, c.z + 1})
	return res
}

func (cm *coordsMap) countFreeSides(c *cube, handleAirPockets bool) int {
	return lo.Reduce(c.neighbours(), func(agg int, item *cube, index int) int {
		if m, set := cm.get(item); set && m == air {
			return agg + 1
		} else if !set {
			if !handleAirPockets {
				// consider it a free side
				return agg + 1
			}
			if cm.expandAir(item) {
				return agg + 1
			}
		}
		return agg
	}, 0)
}

// true if air, false if internal pocket or unknown yet
func (cm *coordsMap) expandAir(c *cube) bool {
	// BFS
	visited := map[cube]struct{}{}
	wave := []*cube{c}
	nextWave := []*cube{}

	// keep track of unknown so we can fill with air/pocket air, once determined
	backQueue := []*cube{}

	var fillMaterial material
	for fill := false; !fill; {
		for _, nextCube := range wave {
			if _, skip := visited[*nextCube]; skip {
				// skip visited (break cycles)
				continue
			}
			visited[*nextCube] = struct{}{}

			if cm.isBeyondLava(nextCube) {
				// set all to air
				fill = true
				fillMaterial = air
				break
			}

			if m, determined := cm.get(nextCube); determined {
				if m == lava {
					// hit lava, skip to next cube
					continue
				}
				// m is pocket air or air - set all to the same
				fill = true
				fillMaterial = m
				break
			}

			// not set so it is unknown, add to back queue
			backQueue = append(backQueue, nextCube)
			nextWave = append(nextWave, nextCube.neighbours()...)
		}

		wave = nextWave
		nextWave = []*cube{}

		if len(wave) == 0 && !fill {
			// no more cubes to expand, must be pocket air
			fill = true
			fillMaterial = pocket_air
		}
	}
	// fill
	for _, nc := range backQueue {
		cm.set(nc, fillMaterial)
	}
	if fillMaterial == air {
		return true
	}

	return false
}

func (cm *coordsMap) isBeyondLava(c *cube) bool {
	if c.x < cm.minX || c.x > cm.maxX ||
		c.y < cm.minY || c.y > cm.maxY ||
		c.z < cm.minZ || c.z > cm.maxZ {
		return true
	}
	return false
}

func processFile(fileName string, handleAirPockets bool) (int, error) {
	cm := NewCoordsMap()

	err := input.ReadFileLines(fileName, func(line string) error {
		ints := input.MustAtoInts(strings.Split(line, ","))
		c := NewCube(ints)
		cm.set(c, lava)
		return nil
	})
	if err != nil {
		return 0, err
	}

	var sides int
	for x, xplane := range cm.m {
		for y, yplane := range xplane {
			for z, m := range yplane {
				if m == lava {
					sides += cm.countFreeSides(&cube{x, y, z}, handleAirPockets)
				}
			}
		}
	}

	return sides, nil
}

func main() {
	res, err := processFile("data/part_one.txt", false)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
	fmt.Println("=================")
	assert.Equals(64, res, "")

	res, err = processFile("data/input.txt", false)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
	fmt.Println("=================")
	assert.Equals(3494, res, "")

	res, err = processFile("data/part_one.txt", true)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
	fmt.Println("=================")
	assert.Equals(58, res, "")

	res, err = processFile("data/input.txt", true)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
	fmt.Println("=================")
	assert.Equals(2062, res, "")
}
