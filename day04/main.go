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

func fullOverlap(a, b secRange) bool {
	if a.s <= b.s {
		// test that a contains b
		return a.e >= b.e
	} else {
		// a.s > b.s
		// test that b contains a
		return a.e <= b.e
	}
}

func procFile(name string) (int, error) {
	var sum int
	var count int
	err := input.ReadFileLines(name, func(line string) error {
		defer func() {
			count++
		}()
		print := func(i interface{}) {
			if count%100 == 0 {
				fmt.Println(i)
			}
		}

		tokens := strings.Split(line, ",")
		if len(tokens) != 2 {
			return errors.New("Wrong line format " + line)
		}

		r1, err := getRange(tokens[0])
		print("---------")
		print(r1)
		if err != nil {
			return err
		}
		r2, err := getRange(tokens[1])
		print(r2)
		if err != nil {
			return err
		}

		if fullOverlap(r1, r2) {
			print("overlap")
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
	res, err := procFile("data/part_one_short.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
	fmt.Println("=================")

	res, err = procFile("data/input.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
	fmt.Println("=================")

}
