package main

import (
	"errors"
	"fmt"
	"strings"

	"kfet.org/aoc_common/assert"
	"kfet.org/aoc_common/input"
)

var (
	rock     = "Rock"
	paper    = "Paper"
	scissors = "Scissors"
)

var theirCode = map[string]string{
	"A": rock,
	"B": paper,
	"C": scissors,
}

var myCode = map[string]string{
	"X": rock,
	"Y": paper,
	"Z": scissors,
}

var scoreMap = map[string]int{
	rock:     1,
	paper:    2,
	scissors: 3,
}

type round struct {
	them string
	me   string
}

var roundMap = map[round]int{
	{rock, rock}:         3,
	{rock, paper}:        6,
	{rock, scissors}:     0,
	{paper, rock}:        0,
	{paper, paper}:       3,
	{paper, scissors}:    6,
	{scissors, rock}:     6,
	{scissors, paper}:    0,
	{scissors, scissors}: 3,
}

func roundScore(them, me string) (int, error) {
	ss, ok := scoreMap[me]
	if !ok {
		return 0, errors.New("Unknown shape score " + me)
	}
	rs, ok := roundMap[round{them, me}]
	if !ok {
		return 0, errors.New("Unknown shapes " + them + ", " + me)
	}

	return ss + rs, nil
}

func myCodeOne(their, mine string) string {
	return myCode[mine]
}

func myCodeTwo(their, mine string) string {
	m := map[round]string{
		{rock, "X"}:     scissors,
		{rock, "Y"}:     rock,
		{rock, "Z"}:     paper,
		{paper, "X"}:    rock,
		{paper, "Y"}:    paper,
		{paper, "Z"}:    scissors,
		{scissors, "X"}: paper,
		{scissors, "Y"}: scissors,
		{scissors, "Z"}: rock,
	}
	return m[round{their, mine}]
}

func strategyScore(name string, codeFn func(string, string) string) (int, error) {
	var score int
	err := input.ReadFileLinesStrings(name, func(tokens []string) error {
		if len(tokens) != 2 {
			return errors.New("Wront number of argumetns " + strings.Join(tokens, " "))
		}
		theirShape := theirCode[tokens[0]]
		myShape := codeFn(theirShape, tokens[1])
		rs, err := roundScore(theirShape, myShape)

		if err != nil {
			return err
		}

		score += rs
		return nil
	})

	if err != nil {
		return 0, err
	}
	return score, nil
}

func main() {
	score, err := strategyScore("data/part_one_small.txt", myCodeOne)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(score)
	fmt.Println("=============")
	assert.Equals(15, score, "")

	score, err = strategyScore("data/input.txt", myCodeOne)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(score)
	fmt.Println("=============")
	assert.Equals(12458, score, "")

	score, err = strategyScore("data/part_one_small.txt", myCodeTwo)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(score)
	fmt.Println("=============")
	assert.Equals(12, score, "")

	score, err = strategyScore("data/input.txt", myCodeTwo)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(score)
	fmt.Println("=============")
	assert.Equals(12683, score, "")
}
