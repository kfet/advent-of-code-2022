package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
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

	switch dir {
	case "L":
		r.knots[0].x -= n
	case "R":
		r.knots[0].x += n
	case "U":
		r.knots[0].y += n
	case "D":
		r.knots[0].y -= n
	}
	r.fixT()
}

var emptyDir knot

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

func (r *rope) fixT() {
	for i := range r.knots {
		if i == 0 {
			// skip head
			continue
		}
		isTail := int64(i) == r.count-1
		prev, cur := &r.knots[i-1], &r.knots[i]
		for dir := prev.pull(cur); dir != emptyDir; dir = prev.pull(cur) {
			cur.x += dir.x
			cur.y += dir.y
			if isTail {
				r.visitTail()
			}
		}
	}
}

func abs[T int | int16 | int32 | int64 | int8](x T) T {
	if x < 0 {
		return -x
	}
	return x
}

func (e *knot) pull(other *knot) knot {
	res := knot{
		x: e.x - other.x,
		y: e.y - other.y,
	}
	aX, aY := abs(res.x), abs(res.y)
	if aX+aY > 2 {
		if aX > 0 {
			res.x = res.x / aX
		}
		if aY > 0 {
			res.y = res.y / aY
		}
	} else {
		squash := func(pi *int64) {
			a := abs(*pi)
			if a > 1 {
				*pi = *pi / a
			} else {
				*pi = 0
			}
		}

		squash(&res.x)
		squash(&res.y)
	}

	return res
}

func (r *rope) runFile(name string) error {
	file, err := os.Open(name)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer file.Close()

	scan := bufio.NewScanner(file)
	for scan.Scan() {
		txt := scan.Text()
		ins := strings.Split(txt, " ")
		if len(ins) != 2 {
			fmt.Printf("Wrong input: %v", txt)
			return os.ErrInvalid
		}

		n, err := strconv.ParseInt(ins[1], 10, 64)
		if err != nil {
			fmt.Printf("Wrong input, can't parse int: %v", txt)
			return os.ErrInvalid
		}

		r.move(ins[0], n)
	}

	if err := scan.Err(); err != nil {
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
	// r := NewRope(2)
	// r.runFile("data/input.txt")
	// fmt.Println(len(r.visited))

	// r = NewRope(2)
	// r.runFile("data/part_one.txt")
	// r.printVisited()
	// fmt.Println(len(r.visited))

	// r = NewRope(10)
	// r.runFile("data/part_two.txt")
	// r.printVisited()
	// fmt.Println(len(r.visited))

	r := NewRope(10)
	r.runFile("data/input.txt")
	fmt.Println(len(r.visited))

	r = NewRope(10)
	r.runFile("data/part_two_short.txt")
	r.printVisited()
	fmt.Println(len(r.visited))
}
