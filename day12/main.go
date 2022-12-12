package main

import (
	"errors"
	"fmt"

	"kfet.org/aoc_common/assert"
	"kfet.org/aoc_common/input"
)

type size struct {
	w, h int
}
type point struct {
	x, y int
}

type questMap struct {
	hm   [][]int
	s, e point
	sz   size
	as   map[point]struct{}
}

func NewQuestMap() *questMap {
	return &questMap{
		as: map[point]struct{}{},
	}
}

func (qm *questMap) readRow(line string, row []int, rowNum int) error {
	for i, r := range line {
		switch {
		case r == 'S':
			row[i] = 0
			qm.s.x = i
			qm.s.y = rowNum
			qm.as[point{x: i, y: rowNum}] = struct{}{}
		case r == 'E':
			row[i] = 'z' - 'a'
			qm.e.x = i
			qm.e.y = rowNum
		case r >= 'a' && r <= 'z':
			row[i] = int(r - 'a')
			if r == 'a' {
				qm.as[point{x: i, y: rowNum}] = struct{}{}
			}
		default:
			return errors.New("Wrong format or row " + fmt.Sprint(rowNum) + ": " + line)
		}
	}

	return nil
}

func (qm *questMap) readMap(file string) error {
	var rowNum int
	err := input.ReadFileLines(file, func(line string) error {
		defer func() {
			rowNum++
		}()

		if len(line) == 0 {
			return nil
		}

		row := make([]int, len(line))
		err := qm.readRow(line, row, rowNum)
		if err != nil {
			return err
		}
		qm.hm = append(qm.hm, row)

		return nil
	})
	if err != nil {
		return err
	}

	qm.sz.w = len(qm.hm[0])
	qm.sz.h = rowNum

	return nil
}

func (p point) isValidMove(n point, qm *questMap) bool {
	if n.x < 0 || n.x >= qm.sz.w ||
		n.y < 0 || n.y >= qm.sz.h {
		return false
	}
	if qm.hm[p.y][p.x]+1 < qm.hm[n.y][n.x] {
		return false
	}
	return true
}

func (qm *questMap) getNeighbours(p point) []point {
	var ns []point
	for _, n := range []point{
		{
			p.x - 1,
			p.y,
		},
		{
			p.x + 1,
			p.y,
		},
		{
			p.x,
			p.y - 1,
		},
		{
			p.x,
			p.y + 1,
		},
	} {
		if p.isValidMove(n, qm) {
			ns = append(ns, n)
		}
	}

	return ns
}

func (qm *questMap) bfsPath(wave map[point]struct{}) (int, bool) {
	var depth int

	visited := map[point]struct{}{}
	nextWave := map[point]struct{}{}

	for {
		for p := range wave {
			if p == qm.e {
				// found
				return depth, true
			}

			// mark as visited
			visited[p] = struct{}{}

			ns := qm.getNeighbours(p)
			for _, n := range ns {
				if _, visited := visited[n]; visited {
					// visited, skip
					continue
				}
				// not visited
				nextWave[n] = struct{}{}
			}
		}

		if len(nextWave) == 0 {
			// not found
			return 0, false
		}

		wave = nextWave
		nextWave = map[point]struct{}{}
		depth++
	}
}

func runQuest(file string, partOne bool) (int, error) {
	qm := *NewQuestMap()

	err := qm.readMap(file)
	if err != nil {
		return 0, err
	}

	var pathLen int
	var foundPath bool
	if partOne {
		// start with S only
		pathLen, foundPath = qm.bfsPath(map[point]struct{}{
			qm.s: {},
		})
	} else {
		// start with all 'a's
		pathLen, foundPath = qm.bfsPath(qm.as)
	}

	if !foundPath {
		return 0, errors.New("path not found")
	}

	return pathLen, nil
}

func main() {
	res, err := runQuest("data/part_one_short.txt", true)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
	fmt.Println("=================")
	assert.Equals(31, res, "")

	res, err = runQuest("data/input.txt", true)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
	fmt.Println("=================")
	assert.Equals(497, res, "")

	res, err = runQuest("data/part_one_short.txt", false)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
	fmt.Println("=================")
	assert.Equals(29, res, "")

	res, err = runQuest("data/input.txt", false)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
	fmt.Println("=================")
	assert.Equals(492, res, "")
}
