package main

import (
	"fmt"
	"math"

	"github.com/samber/lo"
	"kfet.org/aoc_common/assert"
	"kfet.org/aoc_common/input"
)

type digit int
type number []digit

func (n number) toDecimal() int {
	var res int
	l := len(n) - 1
	for i, d := range n {
		p := math.Pow(5, float64(l-i))
		num := int(d) * int(p)
		res += num
	}
	return res
}

func NumberFromString(line string) number {
	return lo.Map([]byte(line), func(item byte, index int) digit {
		switch item {
		case '-':
			return -1
		case '=':
			return -2
		case '2':
			return 2
		case '1':
			return 1
		case '0':
			return 0
		}

		panic(fmt.Sprint("unknown item", item))
	})
}

func NumberFromDecimal(n int) number {
	fiv := n % 5
	if fiv > 2 {
		fiv = fiv - 5
	}

	if n-fiv == 0 {
		return number{digit(fiv)}
	}

	pref := NumberFromDecimal((n - fiv) / 5)
	return append(pref, digit(fiv))
}

func (n number) String() string {
	return lo.Reduce(n, func(agg string, item digit, index int) string {
		var digitChar string
		switch item {
		case -1:
			digitChar = "-"
		case -2:
			digitChar = "="
		case 0, 1, 2:
			digitChar = fmt.Sprint(item)
		default:
			panic(fmt.Sprint("wrong digit ", item))
		}
		return agg + digitChar
	}, "")
}

func processFile(fileName string) (string, error) {
	var resInt int
	err := input.ReadFileLines(fileName, func(line string) error {
		n := NumberFromString(line)
		resInt += n.toDecimal()
		return nil
	})
	if err != nil {
		return "", err
	}
	fmt.Println(resInt)
	return NumberFromDecimal(resInt).String(), nil
}

func main() {

	fmt.Println(NumberFromDecimal(1).String())
	fmt.Println(NumberFromDecimal(2).String())
	fmt.Println(NumberFromDecimal(3).String())
	fmt.Println(NumberFromDecimal(4).String())
	fmt.Println(NumberFromDecimal(5).String())
	fmt.Println(NumberFromDecimal(6).String())
	fmt.Println(NumberFromDecimal(7).String())
	fmt.Println(NumberFromDecimal(8).String())
	fmt.Println(NumberFromDecimal(9).String())
	fmt.Println(NumberFromDecimal(10).String())

	fmt.Println(NumberFromDecimal(2022).String())
	fmt.Println(NumberFromDecimal(12345).String())
	fmt.Println(NumberFromDecimal(314159265).String())
	fmt.Println("=================")

	res, err := processFile("data/part_one.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
	fmt.Println("=================")
	assert.Equals("2=-1=0", res, "")

	res, err = processFile("data/input.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
	fmt.Println("=================")
	assert.Equals("2=1-=02-21===-21=200", res, "")
}
