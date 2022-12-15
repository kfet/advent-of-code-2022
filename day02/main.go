package main

import (
	"errors"
	"fmt"
	"strings"

	"kfet.org/aoc_common/assert"
	"kfet.org/aoc_common/input"
)

type shape string

const (
	rock     = "Rock"
	paper    = "Paper"
	scissors = "Scissors"
)

var theirCode = map[string]shape{
	"A": rock,
	"B": paper,
	"C": scissors,
}

var myCode = map[string]shape{
	"X": rock,
	"Y": paper,
	"Z": scissors,
}

type round struct {
	them shape
	me   shape
}

var (
	shapeScoreMap = map[shape]int{
		rock:     1,
		paper:    2,
		scissors: 3,
	}

	roundScoreMap = map[round]int{
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
)

func roundScore(them, me shape) (int, error) {
	ss, ok := shapeScoreMap[shape(me)]
	if !ok {
		return 0, errors.New("Unknown shape score " + string(me))
	}
	rs, ok := roundScoreMap[round{them, me}]
	if !ok {
		return 0, errors.New("Unknown shapes " + string(them) + ", " + string(me))
	}

	return ss + rs, nil
}

func myCodeQuizOne(their shape, code string) shape {
	return myCode[code]
}

func myCodeQuizTwo(their shape, mine string) shape {
	type codeRound struct {
		their shape
		my    string
	}

	m := map[codeRound]shape{
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
	return m[codeRound{their, mine}]
}

func strategyScore(fileName string, codeFn func(shape, string) shape) (int, error) {
	var score int
	err := input.ReadFileLinesStrings(fileName, func(tokens []string) error {
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
	score, err := strategyScore("data/part_one_small.txt", myCodeQuizOne)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(score)
	fmt.Println("=============")
	assert.Equals(15, score, "")

	score, err = strategyScore("data/input.txt", myCodeQuizOne)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(score)
	fmt.Println("=============")
	assert.Equals(12458, score, "")

	score, err = strategyScore("data/part_one_small.txt", myCodeQuizTwo)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(score)
	fmt.Println("=============")
	assert.Equals(12, score, "")

	score, err = strategyScore("data/input.txt", myCodeQuizTwo)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(score)
	fmt.Println("=============")
	assert.Equals(12683, score, "")
}
