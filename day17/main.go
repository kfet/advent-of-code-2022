package main

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/samber/lo"
	"kfet.org/aoc_common/input"
)

type mask [][]uint8

type rock struct {
	m        *mask
	x, y     int
	moveType int
}

func NewRock(m *mask, x, y int) *rock {
	r := &rock{
		m: m,
		x: x,
		y: y,
	}
	return r
}

const chamberWidth int = 7

type chamber struct {
	w         int
	h         int
	maskStart *big.Int
	rows      mask
}

func NewChamber(width int) *chamber {
	return &chamber{
		w:         width,
		maskStart: big.NewInt(0),
	}
}

func (ch *chamber) String() string {
	bitMap := map[uint8]string{
		0: ".",
		1: "#",
	}
	var sb strings.Builder
	for i := len(ch.rows) - 1; i >= 0; i-- {
		for _, c := range ch.rows[i] {
			sb.WriteString(bitMap[c])
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func (ch *chamber) isNewFloor(row int) bool {
	for _, c := range ch.rows[row] {
		if c == 0 {
			return false
		}
	}
	return true
}

func (ch *chamber) newFloor(row int) {
	ch.rows = ch.rows[row:]
	ch.h = len(ch.rows)
	ch.maskStart.Add(ch.maskStart, big.NewInt(int64(row)))
}

func (ch *chamber) stampRock(r *rock) {
	// add rows as needed
	for ch.h < r.y+1 {
		ch.rows = append(ch.rows, make([]uint8, ch.w))
		ch.h++
	}

	// stamp the rock mask
	var newFloorRow int
	for ry, rr := range *r.m {
		cy := r.y - ry
		for rx, bit := range rr {
			cx := rx + r.x
			if cx < 0 || cx >= ch.w {
				continue
			}
			ch.rows[cy][cx] |= bit
		}

		if ch.isNewFloor(cy) {
			newFloorRow = cy
		}
	}

	if newFloorRow > 0 {
		ch.newFloor(newFloorRow)
	}
}

func (c *chamber) testMove(r *rock, mv move) bool {
	testX := r.x + mv.dx
	testY := r.y + mv.dy
	for ry, rr := range *r.m {
		cy := testY - ry
		if cy < 0 {
			// hit the floor
			return false
		}
		for rx, bit := range rr {
			cx := rx + testX
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
	rockSprites []*mask
	spriteIdx   int

	jets   []move
	curJet int

	ch *chamber
	r  *rock

	rockCount *big.Int
}

const spritesFile string = "data/rocks.txt"

func NewWorld(jetsFile string) *world {
	w := &world{
		spriteIdx: -1,
		ch:        NewChamber(chamberWidth),
		rockCount: big.NewInt(0),
	}

	w.readSprites(spritesFile)
	w.readJets(jetsFile)
	w.nextRock()

	return w
}

func (w *world) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintln("x:", w.r.x, ",y:", w.r.y))
	sb.WriteString(fmt.Sprintln(*w.r.m))
	sb.WriteString(w.ch.String())
	return sb.String()
}

// false if rock stuck
func (w *world) moveRock(r *rock, mv move) bool {
	if w.ch.testMove(r, mv) {
		// can move
		r.x += mv.dx
		r.y += mv.dy
		r.moveType++
		return true
	}
	// can't move

	if mv.dy < 0 {
		// at bottom, stuck
		w.ch.stampRock(r)
		return false
	}

	// next move type
	r.moveType++
	return true
}

func (w *world) step() {
	var nextMove move
	switch w.r.moveType % 2 {
	case 0:
		nextMove = w.jets[w.curJet]
		w.curJet++
		if w.curJet >= len(w.jets) {
			w.curJet = 0
		}

	case 1:
		nextMove = downMove
	}

	if !w.moveRock(w.r, nextMove) {
		// stuck
		w.nextRock()
	}
}

func (w *world) nextSprite() *mask {
	if w.spriteIdx >= len(w.rockSprites)-1 {
		w.spriteIdx = 0
	} else {
		w.spriteIdx++
	}
	return w.rockSprites[w.spriteIdx]
}

var bigOne *big.Int = big.NewInt(1)

func (w *world) nextRock() {
	ns := w.nextSprite()

	x := 2
	y := w.ch.h + len(*ns) + 2

	w.r = NewRock(ns, x, y)
	w.rockCount.Add(w.rockCount, bigOne)
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

func (w *world) readSprites(fileName string) error {
	runeMap := map[rune]uint8{
		'.': 0,
		'#': 1,
	}

	sprite := &mask{}
	err := input.ReadFileLines(fileName, func(line string) error {
		if len(line) == 0 {
			w.rockSprites = append(w.rockSprites, sprite)
			sprite = &mask{}
			return nil
		}

		row := lo.Map([]rune(line), func(r rune, i int) uint8 {
			return runeMap[r]
		})
		(*sprite) = append((*sprite), row)

		return nil
	})
	if err != nil {
		return err
	}
	// append the last rock
	w.rockSprites = append(w.rockSprites, sprite)
	return nil
}

func processFile(fileName string, rockCount *big.Int) (*big.Int, error) {
	var jets string
	err := input.ReadFileLines(fileName, func(line string) error {
		jets = line
		return nil
	})
	if err != nil {
		return nil, err
	}

	w := NewWorld(jets)
	for {
		w.step()
		if w.rockCount.Cmp(rockCount) >= 0 {
			break
		}
	}

	return w.ch.maskStart.Add(w.ch.maskStart, big.NewInt(int64(w.ch.h))), nil
}

func main() {
	res, err := processFile("data/part_one.txt", big.NewInt(2023))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
	fmt.Println("=================")

	res, err = processFile("data/input.txt", big.NewInt(2023))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
	fmt.Println("=================")

	// bigCount, _ := big.NewInt(0).SetString("1000000000000", 10)
	// res, err = processFile("data/part_one.txt", bigCount)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println(res)
	// fmt.Println("=================")
}
