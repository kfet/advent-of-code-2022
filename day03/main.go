package main

import (
	"errors"
	"fmt"

	"kfet.org/aoc_common/assert"
	"kfet.org/aoc_common/input"
)

func runePriority(r rune) int {
	if r >= 'a' && r <= 'z' {
		return int(r-'a') + 1
	}
	if r >= 'A' && r <= 'Z' {
		return int(r-'A') + 27
	}
	return 0
}

func findDuplicate(line string) (rune, error) {
	if len(line)%2 != 0 {
		return 0, errors.New("Odd line length " + line)
	}
	m := map[rune]struct{}{}
	partLen := len(line) / 2
	for idx, r := range line {
		if idx < partLen {
			// store map of runes in first compartment
			m[r] = struct{}{}
		} else {
			// find duplicate rune in second cmpartment
			if _, found := m[r]; !found {
				continue
			}
			// found
			return r, nil
		}
	}
	return 0, errors.New("no duplicate found in line " + line)
}

func partOne(fileName string) (int, error) {
	var sum int
	err := input.ReadFileLines(fileName, func(line string) error {
		r, err := findDuplicate(line)
		if err != nil {
			return err
		}
		sum += runePriority(r)
		return nil
	})
	if err != nil {
		return 0, err
	}
	return sum, nil
}

func partTwo(fileName string) (int, error) {
	var sum int
	var idx int
	var m map[rune]int
	err := input.ReadFileLines(fileName, func(line string) error {
		defer func() {
			idx++
		}()

		switch idx % 3 {
		case 0:
			m = map[rune]int{}
			for _, r := range line {
				m[r] = 0
			}
		case 1:
			for _, r := range line {
				if _, matched := m[r]; matched {
					m[r] = 1
				}
			}
		case 2:
			for _, r := range line {
				if n, matched := m[r]; matched && n == 1 {
					sum += runePriority(r)
					return nil
				}
			}
		}

		return nil
	})
	if err != nil {
		return 0, err
	}
	return sum, nil
}

func main() {
	res, err := partOne("data/part_one_short.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
	fmt.Println("=================")
	assert.Equals(157, res, "")

	res, err = partOne("data/input.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
	fmt.Println("=================")
	assert.Equals(7903, res, "")

	res, err = partTwo("data/part_one_short.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
	fmt.Println("=================")
	assert.Equals(70, res, "")

	res, err = partTwo("data/input.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
	fmt.Println("=================")
	assert.Equals(2548, res, "")
}
