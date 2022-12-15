package main

import (
	"fmt"
	"math"
	"strings"

	"kfet.org/aoc_common/assert"
	"kfet.org/aoc_common/calc"
	"kfet.org/aoc_common/input"
)

type tree struct {
	height                int
	up, down, left, right treeVisibility
	visibilityHook        func(*tree, bool)
}

type treeVisibility struct {
	maxHeight int
	t         *tree
}

type treeMap struct {
	trees     [][]*tree
	visible   map[*tree]struct{}
	invisible map[*tree]struct{}
}

func NewTreeMap() *treeMap {
	tm := &treeMap{
		trees:     make([][]*tree, 0),
		visible:   make(map[*tree]struct{}),
		invisible: map[*tree]struct{}{},
	}
	return tm
}

func (tm *treeMap) String() string {
	var sb strings.Builder
	for _, row := range tm.trees {
		for _, t := range row {
			sb.WriteString(fmt.Sprint(t.height) + " ")
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func (tm *treeMap) registerTree(t *tree) {
	t.visibilityHook = func(t *tree, visible bool) {
		if visible {
			tm.visible[t] = struct{}{}
			delete(tm.invisible, t)
		} else {
			tm.invisible[t] = struct{}{}
			delete(tm.visible, t)
		}
	}
	t.recalcVisibility()
}

func (t *tree) recalcVisibility() {
	if t.height > t.up.maxHeight ||
		t.height > t.left.maxHeight ||
		t.height > t.down.maxHeight ||
		t.height > t.right.maxHeight {
		// visible
		if t.visibilityHook != nil {
			t.visibilityHook(t, true)
		}
		return
	}
	// invisible
	if t.visibilityHook != nil {
		t.visibilityHook(t, false)
	}
}

func (t *tree) recalcVisibilityDir(dir func(*tree) *treeVisibility, back func(*tree) *treeVisibility) {
	v := dir(t)
	if v.maxHeight < dir(v.t).maxHeight ||
		v.maxHeight < v.t.height {
		v.maxHeight = calc.Max(dir(v.t).maxHeight, v.t.height)
		t.recalcVisibility()

		if bt := back(t); bt.t != nil {
			bt.t.recalcVisibilityDir(dir, back)
		}
	}
}

func (t *tree) recalcVisibilityDown() {
	t.recalcVisibilityDir(
		func(t *tree) *treeVisibility { return &t.down },
		func(t *tree) *treeVisibility { return &t.up },
	)
}

func (t *tree) recalcVisibilityRight() {
	t.recalcVisibilityDir(
		func(t *tree) *treeVisibility { return &t.right },
		func(t *tree) *treeVisibility { return &t.left },
	)
}

func NewTree(up, left *tree, height int) *tree {
	t := &tree{height: height,
		up:    treeVisibility{-1, nil},
		left:  treeVisibility{-1, nil},
		down:  treeVisibility{-1, nil},
		right: treeVisibility{-1, nil},
	}

	// up/down
	if up != nil {
		// up visiblity based on existing calcs to the top
		t.up = treeVisibility{
			t:         up,
			maxHeight: calc.Max(up.up.maxHeight, up.height),
		}
		// hook down link and recalc column down visibility
		up.down.t = t
		up.recalcVisibilityDown()
	}

	// left/right
	if left != nil {
		// left visiblity based on existing calcs to the left
		t.left = treeVisibility{
			t:         left,
			maxHeight: calc.Max(left.left.maxHeight, left.height),
		}
		// hook right link and recalc row right visibility
		left.right.t = t
		left.recalcVisibilityRight()
	}

	return t
}

func processFile(name string, partOne bool) (int, error) {

	tm := NewTreeMap()

	var row int
	err := input.ReadFileLines(name, func(line string) error {
		// add new row
		tm.trees = append(tm.trees, make([]*tree, 0))

		// add each tree to the row
		for col, r := range line {
			var up *tree
			if row > 0 {
				up = tm.trees[row-1][col]
			}
			var left *tree
			if col > 0 {
				left = tm.trees[row][col-1]
			}
			h := int(r - '0')

			t := NewTree(up, left, h)
			tm.trees[row] = append(tm.trees[row], t)
			tm.registerTree(t)
		}
		row++
		return nil
	})
	if err != nil {
		return 0, err
	}

	if partOne {
		return len(tm.visible), nil
	}

	// part two
	scoreInDirection := func(t *tree, nextTree func(*tree) *tree) int {
		nt := nextTree(t)
		var dirScore int
		for nt != nil {
			dirScore++
			if nt.height >= t.height {
				// tree hides the view
				break
			}
			nt = nextTree(nt)
		}
		return dirScore
	}

	score := func(t *tree) int {
		u := scoreInDirection(t, func(t *tree) *tree { return t.up.t })
		l := scoreInDirection(t, func(t *tree) *tree { return t.left.t })
		d := scoreInDirection(t, func(t *tree) *tree { return t.down.t })
		r := scoreInDirection(t, func(t *tree) *tree { return t.right.t })

		return u * l * d * r
	}

	maxScore := math.MinInt
	for _, r := range tm.trees {
		for _, t := range r {
			s := score(t)
			if s > maxScore {
				maxScore = s
			}
		}
	}

	return maxScore, nil
}

func main() {
	res, err := processFile("data/part_one.txt", true)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
	fmt.Println("=================")
	assert.Equals(21, res, "")

	res, err = processFile("data/input.txt", true)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
	fmt.Println("=================")
	assert.Equals(1708, res, "")

	res, err = processFile("data/part_one.txt", false)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
	fmt.Println("=================")

	res, err = processFile("data/input.txt", false)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
	fmt.Println("=================")
}
