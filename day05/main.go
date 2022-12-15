package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"kfet.org/aoc_common/assert"
	"kfet.org/aoc_common/input"
)

type stacks []stack
type stack []*crate
type crate rune

func NewCrate(r rune) *crate {
	c := new(crate)
	*c = crate(r)
	return c
}

func (s stack) String() string {
	var sb strings.Builder
	for _, c := range s {
		sb.WriteRune(rune(*c))
		sb.WriteRune(',')
	}
	return sb.String()
}

func (s stacks) String() string {
	var sb strings.Builder
	for _, s := range s {
		sb.WriteString(s.String())
		sb.WriteRune('\n')
	}
	return sb.String()
}

func moveOne(s *stacks, count, from, to int) {
	// move the crates one by one
	for i := 0; i < count; i++ {
		// pop from 'from'
		c := (*s)[from][0]
		(*s)[from] = (*s)[from][1:]
		// prepend to 'to'
		(*s)[to] = append([]*crate{NewCrate(rune(*c))}, (*s)[to]...)
	}
}

func moveTwo(s *stacks, count, from, to int) {
	// move the crates all at once
	// pop 'count' crates
	c := (*s)[from][0:count]
	fs := make(stack, 0)
	fs = append(fs, (*s)[from][count:]...)
	(*s)[from] = fs

	// prepend them to 'to'
	fs = make(stack, 0)
	fs = append(fs, c...)
	fs = append(fs, (*s)[to]...)
	(*s)[to] = fs
}

func (s *stacks) moveCrates(scan *bufio.Scanner, moveFunc func(*stacks, int, int, int)) error {

	for scan.Scan() {
		line := scan.Text()
		tokens := regexp.MustCompile(`move (\d+) from (\d+) to (\d+)`).FindStringSubmatch(line)
		if tokens == nil {
			return errors.New("wrong move line format " + line)
		}
		nums := input.MustAtoInts(tokens[1:])
		moveFunc(s, nums[0], nums[1]-1, nums[2]-1)
	}
	if err := scan.Err(); err != nil {
		return err
	}

	return nil
}

func (s *stacks) readCrates(line string) error {
	nStacks := (len(line) + 1) / 4
	for len(*s) < nStacks {
		*s = append(*s, make(stack, 0))
	}

	for i := range *s {
		r := rune(line[i*4+1])
		if r != ' ' {
			(*s)[i] = append((*s)[i], NewCrate(r))
		}
	}

	return nil
}

func readStacks(scan *bufio.Scanner) (stacks, error) {
	s := make(stacks, 0)

	for scan.Scan() {
		var line string
		if line = scan.Text(); len(line) == 0 {
			// End of stacks diagram
			return s, nil
		}
		// len(line) > 0
		err := s.readCrates(line)
		if err != nil {
			return nil, nil
		}
	}
	if err := scan.Err(); err != nil {
		return nil, err
	}

	return nil, errors.New("wrong file format, no end of crate stacks")
}

func processFile(name string, moveFunc func(*stacks, int, int, int)) (string, error) {
	file, err := os.Open(name)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scan := bufio.NewScanner(file)

	stacks, err := readStacks(scan)
	if err != nil {
		return "", err
	}

	err = stacks.moveCrates(scan, moveFunc)
	if err != nil {
		return "", err
	}

	var sb strings.Builder
	for _, st := range stacks {
		r := rune(*st[0])
		sb.WriteRune(r)
	}

	return sb.String(), nil
}

func main() {
	res, err := processFile("data/part_one.txt", moveOne)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
	fmt.Println("=================")
	assert.Equals("CMZ", res, "")

	res, err = processFile("data/input.txt", moveOne)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
	fmt.Println("=================")
	assert.Equals("LJSVLTWQM", res, "")

	res, err = processFile("data/part_one.txt", moveTwo)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
	fmt.Println("=================")
	assert.Equals("MCD", res, "")

	res, err = processFile("data/input.txt", moveTwo)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
	fmt.Println("=================")
	assert.Equals("BRQWDBBJM", res, "")
}
