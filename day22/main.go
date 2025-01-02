package main

import (
	"context"
	"fmt"

	"kfet.org/aoc_common/input"
)

type fieldMap struct {
}

func NewFieldMap() *fieldMap {
	return &fieldMap{}
}

func (fm *fieldMap) readMap(row int, line string) bool {
	if len(line) == 0 {
		// end of map
		return false
	}

	// TODO
	return true
}

type direction uint8

const (
	right = iota
	down
	left
	up
)

type move struct {
	dir   direction
	steps int
}

type path []move

// 10R5L5R10L4R5L5
// assume first letter is "R"
var dirRunes = map[rune]direction{
	'R': right,
	// 'L':
}

func readPath(line string) *path {
	// TODO
	for i, r := range line {
		_, _ = i, r
	}
	return nil
}

func processFile(fileName string) (int, error) {
	fm := NewFieldMap()
	var row int
	isMap := true
	err := input.ReadFileLines(fileName, func(line string) error {
		isMap = fm.readMap(row, line)
		if isMap {
			row++
			return nil
		}
		// !isMap
		if len(line) > 0 {
			readPath(line)
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return 0, nil
}

func main() {
	ctx := context.Background()
	c, cancel := context.WithCancel(ctx)
	cancel()
	vc := context.WithValue(c, "asfd-key", "asdf,val")
	_ = vc

	res, err := processFile("data/part_one.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
	fmt.Println("=================")

	res, err = processFile("data/input.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
	fmt.Println("=================")
}
