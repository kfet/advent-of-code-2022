package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"strconv"
)

// Adopted IntHeap example from container/heap
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

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
		return elfs, 0, err
	}

	scan := bufio.NewScanner(file)
	var elfCalories int
	for scan.Scan() {
		txt := scan.Text()
		if len(txt) == 0 {
			// new line
			heap.Push(&elfs, elfCalories)
			elfCalories = 0
			continue
		}

		cals, err := strconv.ParseInt(txt, 10, 64)
		if err != nil {
			fmt.Printf("Failed parsing number %v", txt)
			fmt.Println(err)
			return elfs, 0, err
		}
		elfCalories += int(cals)
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
	fmt.Printf("Part one example: %v\n", cals)

	_, cals, err = maxElfsCalories("data/input.txt", 1)
	if err != nil {
		fmt.Print(err)
		return
	}
	fmt.Printf("Part one all input: %v\n", cals)

	_, cals, err = maxElfsCalories("data/part_one.txt", 3)
	if err != nil {
		fmt.Print(err)
		return
	}
	fmt.Printf("Part two example: %v\n", cals)

	_, cals, err = maxElfsCalories("data/input.txt", 3)
	if err != nil {
		fmt.Print(err)
		return
	}
	fmt.Printf("Part two all input: %v\n", cals)
}
