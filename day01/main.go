package main

import (
	"container/heap"
	"fmt"

	"kfet.org/aoc_common/assert"
	"kfet.org/aoc_common/input"
)

type elfMaxHeap []int

func (e elfMaxHeap) Len() int           { return len(e) }
func (e elfMaxHeap) Less(i, j int) bool { return e[i] > e[j] }
func (e elfMaxHeap) Swap(i, j int)      { e[i], e[j] = e[j], e[i] }

func (e *elfMaxHeap) Push(x any) {
	*e = append(*e, x.(int))
}

func (e *elfMaxHeap) Pop() any {
	old := *e
	n := len(old)
	x := old[n-1]
	*e = old[0 : n-1]
	return x
}

func maxElfsCalories(fileName string, maxElfs int) (elfMaxHeap, int, error) {

	var elfs elfMaxHeap
	heap.Init(&elfs)

	var elfCalories int

	err := input.ReadFileLines(fileName, func(line string) error {
		if len(line) == 0 {
			// new line
			heap.Push(&elfs, elfCalories)
			elfCalories = 0
			return nil
		}

		elfCalories += input.MustAtoi(line)
		return nil
	})

	if err != nil {
		return nil, 0, err
	}

	var maxCalories int
	for i := 0; i < maxElfs; i++ {
		maxCalories += heap.Pop(&elfs).(int)
	}
	return elfs, maxCalories, nil
}

func main() {
	_, cals, err := maxElfsCalories("data/part_one.txt", 1)
	if err != nil {
		fmt.Print(err)
		return
	}

	assert.Equals(24000, cals, "")
	fmt.Printf("Part one example: %v\n", cals)

	_, cals, err = maxElfsCalories("data/input.txt", 1)
	if err != nil {
		fmt.Print(err)
		return
	}
	assert.Equals(66616, cals, "")
	fmt.Printf("Part one all input: %v\n", cals)

	_, cals, err = maxElfsCalories("data/part_one.txt", 3)
	if err != nil {
		fmt.Print(err)
		return
	}
	assert.Equals(41000, cals, "")
	fmt.Printf("Part two example: %v\n", cals)

	_, cals, err = maxElfsCalories("data/input.txt", 3)
	if err != nil {
		fmt.Print(err)
		return
	}
	assert.Equals(199172, cals, "")
	fmt.Printf("Part two all input: %v\n", cals)
}
