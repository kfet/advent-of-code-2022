package main

import (
	"fmt"

	"kfet.org/aoc_common/input"
)

func processFile(fileName string) (int, error) {
	err := input.ReadFileLines(fileName, func(line string) error {
		return nil
	})
	if err != nil {
		return 0, err
	}
	return 0, nil
}

func main() {
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
