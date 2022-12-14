package main

import (
	"errors"
	"fmt"

	"kfet.org/aoc_common/assert"
	"kfet.org/aoc_common/input"
)

func isMarker(slice string) bool {
	m := map[rune]struct{}{}
	for _, r := range slice {
		if _, found := m[r]; found {
			return false
		}
		m[r] = struct{}{}
	}
	return true
}

func findMarker(line string, num int) (int, error) {
	var i int
	for i+num-1 < len(line) {
		if isMarker(line[i : i+num]) {
			return i + num - 1, nil
		}
		i++
	}
	return 0, errors.New("marker not found")
}

func processFile(name string, num int) (int, error) {
	var idx int
	err := input.ReadFileLines(name, func(line string) error {
		i, err := findMarker(line, num)
		if err != nil {
			return err
		}
		idx = i
		return nil
	})
	if err != nil {
		return -1, err
	}
	return idx + 1, nil
}

func main() {
	res, err := processFile("data/part_one.txt", 4)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
	fmt.Println("=================")
	assert.Equals(7, res, "")

	res, err = findMarker("bvwbjplbgvbhsrlpgdmjqwftvncz", 4)
	fmt.Println(res+1, err)
	fmt.Println("=================")
	assert.Equals(5, res+1, "")

	res, err = findMarker("nppdvjthqldpwncqszvftbrmjlhg", 4)
	fmt.Println(res+1, err)
	fmt.Println("=================")
	assert.Equals(6, res+1, "")

	res, err = findMarker("nznrnfrfntjfmvfwmzdfjlvtqnbhcprsg", 4)
	fmt.Println(res+1, err)
	fmt.Println("=================")
	assert.Equals(10, res+1, "")

	res, err = findMarker("zcfzfwzzqfrljwzlrfnpqdbhtmscgvjw", 4)
	fmt.Println(res+1, err)
	fmt.Println("=================")
	assert.Equals(11, res+1, "")

	res, err = processFile("data/input.txt", 4)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
	fmt.Println("=================")
	assert.Equals(1850, res, "")

	res, err = processFile("data/part_one.txt", 14)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
	fmt.Println("=================")
	assert.Equals(19, res, "")

	res, err = findMarker("bvwbjplbgvbhsrlpgdmjqwftvncz", 14)
	fmt.Println(res+1, err)
	fmt.Println("=================")
	assert.Equals(23, res+1, "")

	res, err = findMarker("nppdvjthqldpwncqszvftbrmjlhg", 14)
	fmt.Println(res+1, err)
	fmt.Println("=================")
	assert.Equals(23, res+1, "")

	res, err = findMarker("nznrnfrfntjfmvfwmzdfjlvtqnbhcprsg", 14)
	fmt.Println(res+1, err)
	fmt.Println("=================")
	assert.Equals(29, res+1, "")

	res, err = findMarker("zcfzfwzzqfrljwzlrfnpqdbhtmscgvjw", 14)
	fmt.Println(res+1, err)
	fmt.Println("=================")
	assert.Equals(26, res+1, "")

	res, err = processFile("data/input.txt", 14)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
	fmt.Println("=================")
	assert.Equals(2823, res, "")
}
