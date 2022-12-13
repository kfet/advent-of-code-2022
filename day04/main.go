package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"kfet.org/aoc_common/input"
)

type secRange struct {
	s, e int
}

func getRange(s string) (secRange, error) {
	tokens := strings.Split(s, "-")
	if len(tokens) != 2 {
		return secRange{}, errors.New("Wrong range format " + s)
	}
	a, err := strconv.Atoi(tokens[0])
	if err != nil {
		return secRange{}, err
	}
	b, err := strconv.Atoi(tokens[1])
	if err != nil {
		return secRange{}, err
	}
	return secRange{s: a, e: b}, nil
}

func anyOverlap(a, b secRange) bool {
	if a.s == b.s {
		return true
	}
	if a.s < b.s {
		return a.e >= b.s
	}
	// a.s > b.s
	return a.s <= b.e
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
		tokens := strings.Split(line, ",")
		if len(tokens) != 2 {
			return errors.New("Wrong line format " + line)
		}

		r1, err := getRange(tokens[0])
		if err != nil {
			return err
		}
		r2, err := getRange(tokens[1])
		if err != nil {
			return err
		}

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
