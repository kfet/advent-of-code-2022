package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"kfet.org/aoc_common/assert"
)

type item struct {
	worry int64
}

var validOps = map[string]struct{}{
	"+": {},
	"*": {},
	"^": {},
}

type op struct {
	name string
	val  int64
}

type test struct {
	divBy   int64
	posDest int
	negDest int
}

type monkey struct {
	activity int
	items    []item
	op       op
	test     test
}

var (
	modulo  int64
	monkeys []*monkey
)

func readMonkey(scan *bufio.Scanner) (*monkey, error) {

	var m monkey

	// Monkey %d:
	mNo, err := matchLine(scan, "Monkey ")
	if err != nil {
		return nil, err
	}
	if len(mNo) == 0 {
		// EOF
		return nil, nil
	}

	// Items
	itemsStr, err := matchLine(scan, "  Starting items: ")
	if err != nil {
		return nil, err
	}
	tokens := strings.Split(itemsStr, ", ")
	for _, itemStr := range tokens {
		worryLevel, err := strconv.ParseInt(itemStr, 10, 32)
		if err != nil {
			return nil, err
		}

		m.items = append(m.items, item{
			worry: worryLevel,
		})
	}

	// operation
	opStr, err := matchLine(scan, "  Operation: new = old ")
	if err != nil {
		return nil, err
	}
	if opStr == "* old" {
		m.op = op{
			name: "^",
		}
	} else if opStr == "+ old" {
		m.op = op{
			name: "*",
			val:  2,
		}
	} else {
		tokens = strings.Split(opStr, " ")
		if len(tokens) != 2 {
			return nil, errors.New("Wrong op format " + strings.Join(tokens, " "))
		}
		if _, ok := validOps[tokens[0]]; !ok {
			return nil, errors.New("Wrong op name " + strings.Join(tokens, " "))
		}
		val, err := strconv.ParseInt(tokens[1], 10, 32)
		if err != nil {
			return nil, err
		}
		m.op = op{
			name: tokens[0],
			val:  val,
		}
	}

	// test
	num, err := matchLine(scan, "  Test: divisible by ")
	if err != nil {
		return nil, err
	}
	divBy, err := strconv.ParseInt(num, 10, 32)
	if err != nil {
		return nil, err
	}
	m.test.divBy = divBy

	num, err = matchLine(scan, "    If true: throw to monkey ")
	if err != nil {
		return nil, err
	}
	posDest, err := strconv.ParseInt(num, 10, 32)
	if err != nil {
		return nil, err
	}
	m.test.posDest = int(posDest)

	num, err = matchLine(scan, "    If false: throw to monkey ")
	if err != nil {
		return nil, err
	}
	negDest, err := strconv.ParseInt(num, 10, 32)
	if err != nil {
		return nil, err
	}
	m.test.negDest = int(negDest)

	return &m, nil
}

func matchLine(scan *bufio.Scanner, match string) (string, error) {
	var line string
	// Just ignore all empty lines
	for len(line) == 0 {
		if !scan.Scan() {
			return "", nil
		}
		line = scan.Text()
	}

	if len(line) < len(match) {
		return "", errors.New("Wrong lengh " + line + ", for match " + match)
	}
	if line[0:len(match)] != match {
		return "", errors.New("Wrong format " + line)
	}

	return line[len(match):], nil
}

func processFile(name string, worryFactor int, rounds int) (int, error) {
	file, err := os.Open(name)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scan := bufio.NewScanner(file)

	monkeys = []*monkey{}
	modulo = 1

	m, err := readMonkey(scan)
	if err != nil {
		return 0, err
	}
	for m != nil {
		monkeys = append(monkeys, m)
		modulo *= m.test.divBy

		m, err = readMonkey(scan)
		if err != nil {
			return 0, err
		}
	}

	if len(monkeys) == 0 {
		return 0, nil
	}

	for i := 0; i < rounds; i++ {
		for _, m := range monkeys {
			monkeyRound(m, worryFactor)
		}
	}

	if len(monkeys) == 1 {
		return monkeys[0].activity, nil
	}
	// len(monkeys) >= 2

	sort.Slice(monkeys, func(i, j int) bool {
		return monkeys[i].activity > monkeys[j].activity
	})

	return monkeys[0].activity * monkeys[1].activity, nil
}

func monkeyRound(m *monkey, worryFactor int) {
	for len(m.items) > 0 {
		m.inspectItem(worryFactor)
		m.handleItem()
	}
}

func (m *monkey) inspectItem(worryFactor int) {
	m.activity++

	item := &m.items[0]
	switch m.op.name {
	case "+":
		item.worry += m.op.val
	case "*":
		item.worry *= m.op.val
	case "^":
		// old * old
		item.worry *= item.worry
	}

	if worryFactor != 1 {
		item.worry /= int64(worryFactor)
	}

	item.worry %= modulo
}

func (m *monkey) handleItem() {
	item := m.items[0]
	if item.worry%m.test.divBy == 0 {
		// move to pos dest
		monkeys[m.test.posDest].items = append(monkeys[m.test.posDest].items, item)
	} else {
		// move to neg dest
		monkeys[m.test.negDest].items = append(monkeys[m.test.negDest].items, item)
	}
	m.items = m.items[1:]
}

func printMonkeys() {
	for _, m := range monkeys {
		fmt.Println(*m)
	}
}

func main() {
	res, err := processFile("data/part_one_short.txt", 3, 20)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
	fmt.Println("=================")
	assert.Equals(10605, res, "")

	res, err = processFile("data/input.txt", 3, 20)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
	fmt.Println("=================")
	assert.Equals(110220, res, "")

	res, err = processFile("data/part_one_short.txt", 1, 10000)
	if err != nil {
		fmt.Println(err)
		return
	}
	printMonkeys()
	fmt.Println(res)
	fmt.Println("=================")
	assert.Equals(2713310158, res, "")

	res, err = processFile("data/input.txt", 1, 10000)
	if err != nil {
		fmt.Println(err)
		return
	}
	printMonkeys()
	fmt.Println(res)
	fmt.Println("=================")
	assert.Equals(19457438264, res, "")
}
