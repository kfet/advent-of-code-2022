package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"kfet.org/aoc_common/assert"
	"kfet.org/aoc_common/calc"
	"kfet.org/aoc_common/input"
)

type knot struct {
	x int64
	y int64
}

type point knot
type size struct {
	bottomLeft point
	topRight   point
}

type rope struct {
	knots   []knot
	count   int64
	visited map[knot]struct{}
	size    size
}

func NewRope(knotCount int64) *rope {
	r := &rope{}
	r.count = knotCount
	r.knots = make([]knot, knotCount)
	r.visited = map[knot]struct{}{
		r.knots[r.count-1]: {},
	}
	return r
}

func (r *rope) move(dir string, n int64) {

	for i := int64(0); i < n; i++ {
		switch dir {
		case "L":
			r.knots[0].x--
		case "R":
			r.knots[0].x++
		case "U":
			r.knots[0].y++
		case "D":
			r.knots[0].y--
		}
		r.pullRope()
	}
}

var noOp knot

func (r *rope) visitTail() {
	pad := int64(2)
	tail := &r.knots[r.count-1]
	r.visited[*tail] = struct{}{}
	if tail.x-pad < r.size.bottomLeft.x {
		r.size.bottomLeft.x = tail.x - pad
	}
	if tail.x+pad > r.size.topRight.x {
		r.size.topRight.x = tail.x + pad
	}
	if tail.y-pad < r.size.bottomLeft.y {
		r.size.bottomLeft.y = tail.y - pad
	}
	if tail.y+pad > r.size.topRight.y {
		r.size.topRight.y = tail.y + pad
	}
}

func (r *rope) pullRope() {
	for i := 1; i < len(r.knots); i++ {
		isTail := int64(i) == r.count-1
		prev, cur := &r.knots[i-1], &r.knots[i]
		for dir := prev.pull(cur); dir != noOp; dir = prev.pull(cur) {
			cur.x += dir.x
			cur.y += dir.y
			if isTail {
				r.visitTail()
			}
		}
	}
}

func (k *knot) pull(other *knot) knot {
	res := knot{
		x: k.x - other.x,
		y: k.y - other.y,
	}

	aX, aY := calc.Abs(res.x), calc.Abs(res.y)
	if aX+aY > 2 {
		// diagonal moves allowed
		if aX > 0 {
			res.x = res.x / aX
		}
		if aY > 0 {
			res.y = res.y / aY
		}
	} else {
		squash := func(pi *int64) {
			a := calc.Abs(*pi)
			if a > 1 {
				// pull in this direction
				*pi = *pi / a
			} else {
				// adjacent, do not in this direction
				*pi = 0
			}
		}

		squash(&res.x)
		squash(&res.y)
	}

	return res
}

func (r *rope) runFile(name string) error {

	err := input.ReadFileLines(name, func(line string) error {
		ins := strings.Split(line, " ")
		if len(ins) != 2 {
			fmt.Printf("Wrong input: %v", line)
			return os.ErrInvalid
		}

		n, err := strconv.ParseInt(ins[1], 10, 64)
		if err != nil {
			fmt.Printf("Wrong input, can't parse int: %v", line)
			return os.ErrInvalid
		}

		r.move(ins[0], n)
		return nil
	})

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (r *rope) printVisited() {
	for i := r.size.topRight.y; i >= r.size.bottomLeft.y; i-- {
		var line strings.Builder
		for j := r.size.bottomLeft.x; j <= r.size.topRight.x; j++ {
			if _, exists := r.visited[knot{x: int64(j), y: int64(i)}]; exists {
				line.WriteString("#")
			} else {
				line.WriteString(".")
			}
		}
		fmt.Println(line.String())
	}
}

func main() {
	r := NewRope(2)
	r.runFile("data/input.txt")

	assert.Equals(6057, len(r.visited), "")
	fmt.Println(len(r.visited))

	r = NewRope(2)
	r.runFile("data/part_one.txt")
	r.printVisited()
	assert.Equals(13, len(r.visited), "")
	fmt.Println(len(r.visited))

	r = NewRope(10)
	r.runFile("data/part_two.txt")
	r.printVisited()
	assert.Equals(36, len(r.visited), "")
	fmt.Println(len(r.visited))

	r = NewRope(10)
	r.runFile("data/input.txt")
	assert.Equals(2514, len(r.visited), "")
	fmt.Println(len(r.visited))

	r = NewRope(10)
	r.runFile("data/part_two_short.txt")
	r.printVisited()
	assert.Equals(1, len(r.visited), "")
	fmt.Println(len(r.visited))
}
