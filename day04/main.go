package main

import (
	"errors"
	"fmt"
	"regexp"

	"kfet.org/aoc_common/input"
)

type secRange struct {
	s, e int
}

func NewSecRange(i []int) *secRange {
	if i[1] < i[0] {
		return &secRange{s: i[1], e: i[0]}
	}
	return &secRange{s: i[0], e: i[1]}
}

func anyOverlap(a, b secRange) bool {
	if a.e < b.s || b.e < a.s {
		return false
	}
	return true
}

func fullOverlap(a, b secRange) bool {
	if a.s == b.s {
		return true
	}
	if a.s < b.s {
		// test that a contains b
		return a.e >= b.e
	}
	// a.s > b.s
	// test that b contains a
	return a.e <= b.e

}

func processFile(name string, testFunc func(secRange, secRange) bool) (int, error) {
	var sum int

	err := input.ReadFileLines(name, func(line string) error {
		re := regexp.MustCompile(`^(\d+)-(\d+),(\d+)-(\d+)$`)
		tokens := re.FindStringSubmatch(line)
		if tokens == nil {
			return errors.New("Wrong line format " + line)
		}

		r1 := *NewSecRange(input.MustAtoInts(tokens[1:3]))
		r2 := *NewSecRange(input.MustAtoInts(tokens[3:5]))

		if testFunc(r1, r2) {
			sum++
		}
		return nil
	})
	if err != nil {
		return 0, err
	}

	return sum, nil
}

func main() {
	res, err := processFile("data/part_one_short.txt", fullOverlap)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
	fmt.Println("=================")

	res, err = processFile("data/input.txt", fullOverlap)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
	fmt.Println("=================")

	res, err = processFile("data/part_one_short.txt", anyOverlap)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
	fmt.Println("=================")

	res, err = processFile("data/input.txt", anyOverlap)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
	fmt.Println("=================")
}
