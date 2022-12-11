package main

import (
	"fmt"

	"kfet.org/aoc_common/input"
)

func procFile(name string) (int, error) {
	err := input.ReadFileLinesStrings(name, func(tokens []string) error {
		return nil
	})
	if err != nil {
		return 0, err
	}
	return 0, nil
}

func main() {
	res, err := procFile("data/part_one_short.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
	fmt.Println("=================")

	res, err = procFile("data/part_one.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
	fmt.Println("=================")

	res, err = procFile("data/part_two.txt")
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
