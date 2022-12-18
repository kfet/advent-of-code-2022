package main

import (
	"fmt"
	"strings"

	"github.com/samber/lo"
	"kfet.org/aoc_common/input"
)

type mask [][]uint8

type rock mask

const chamberWidth int = 7

type chamber struct {
	w    int
	h    int
	rows mask
}

func NewChamber(width int) *chamber {
	return &chamber{
		w: width,
	}
}

func (c *chamber) stampRock(r *rock, x, y int) {
	fmt.Println("Stamp rock")
	fmt.Println(x, y)
	fmt.Println(*r)
	fmt.Println()

	// add rows as needed
	for c.h < y+1 {
		c.rows = append(c.rows, make([]uint8, c.w))
		c.h++
	}

	// stamp the rock mask
	for ry, rr := range *r {
		cy := y - ry
		for rx, bit := range rr {
			cx := rx + x
			if cx < 0 || cx >= c.w {
				continue
			}
			c.rows[cy][cx] |= bit
		}
	}
}

func (c *chamber) testMove(r *rock, x, y int) bool {
	for ry, rr := range *r {
		cy := y - ry
		if cy < 0 {
			// hit the floor
			return false
		}
		for rx, bit := range rr {
			cx := rx + x
			if cx < 0 || cx >= c.w {
				// attempt to move outside of chamber
				return false
			}
			if cy >= c.h {
				// above chamber top
				continue
			}
			if c.rows[cy][cx]&bit != 0 {
				// hit a rock
				return false
			}
		}
	}

	return true
}

var downMove move = move{0, -1}

type move struct {
	dx, dy int
}

type world struct {
	rocks    []*rock
	rockIdx  int
	x, y     int
	rockMove int

	jets   []move
	curJet int

	ch *chamber
}

const rocksFile string = "data/rocks.txt"

func NewWorld(jetsFile string) *world {
	w := &world{
		ch:    NewChamber(chamberWidth),
		rocks: []*rock{},

		rockIdx: -1,
	}

	w.readRocks(rocksFile)
	w.readJets(jetsFile)
	w.nextRock()

	return w
}

func (w *world) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintln("x:", w.x, ",y:", w.y))
	sb.WriteRune('\n')
	for i := len(w.ch.rows) - 1; i >= 0; i-- {
		r := w.ch.rows[i]
		sb.WriteString(fmt.Sprintln(r))
	}
	sb.WriteRune('\n')
	return sb.String()
}

func (w *world) step() {

	var nextMove move
	switch w.rockMove % 2 {
	case 0:
		fmt.Println("Jet move")
		nextMove = w.jets[w.curJet]
		w.curJet++
		if w.curJet >= len(w.jets) {
			w.curJet = 0
		}

	case 1:
		fmt.Println("Down move")
		nextMove = downMove
	}

	if w.ch.testMove(w.rocks[w.rockIdx], w.x+nextMove.dx, w.y+nextMove.dy) {
		// can move
		w.x += nextMove.dx
		w.y += nextMove.dy
	} else {
		// hit rock
		if w.rockMove == 1 {
			// ... at bottom
			w.ch.stampRock(w.rocks[w.rockIdx], w.x, w.y)
			w.nextRock()
			w.rockMove = -1
		}
	}

	// next move
	w.rockMove++
}

func (w *world) nextRock() {
	if w.rockIdx >= len(w.rocks)-1 {
		w.rockIdx = 0
	} else {
		w.rockIdx++
	}

	w.x = 2

	rh := len(*w.rocks[w.rockIdx])
	w.y = w.ch.h + rh + 2

	w.rockMove = 0
}

func (w *world) readJets(line string) {
	moveMap := map[rune]move{
		'<': {-1, 0},
		'>': {+1, 0},
	}
	w.jets = lo.Map([]rune(line), func(item rune, index int) move {
		return moveMap[item]
	})
}

func (w *world) readRocks(fileName string) error {
	runeMap := map[rune]uint8{
		'.': 0,
		'#': 1,
	}

	r := &rock{}
	err := input.ReadFileLines(fileName, func(line string) error {
		if len(line) == 0 {
			w.rocks = append(w.rocks, r)
			r = &rock{}
			return nil
		}

		row := lo.Map([]rune(line), func(r rune, i int) uint8 {
			return runeMap[r]
		})
		(*r) = append((*r), row)

		return nil
	})
	if err != nil {
		return err
	}
	// append the last rock
	w.rocks = append(w.rocks, r)
	return nil
}

func processFile(fileName string) (int, error) {
	var jets string
	err := input.ReadFileLines(fileName, func(line string) error {
		jets = line
		return nil
	})
	if err != nil {
		return 0, err
	}

	w := NewWorld(jets)

	fmt.Println(w.jets)
	for _, r := range w.rocks {
		fmt.Println(r)
	}
	fmt.Println()

	for i := 0; i < 10; i++ {
		fmt.Println(w.String())
		w.step()
	}
	fmt.Println(w.String())

	return 0, nil
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
