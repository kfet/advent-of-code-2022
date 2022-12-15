package main

import (
	"bufio"
	"container/heap"
	"errors"
	"fmt"
	"os"
	"strings"

	"kfet.org/aoc_common/assert"
	"kfet.org/aoc_common/input"
)

type item struct {
	isValue bool // true if value, false if list
	value   int
	list    items
}

func (i item) String() string {
	if !i.isValue {
		return i.list.String()
	}
	return fmt.Sprintf("%v,", i.value)
}

type items []*item

func (its items) String() string {
	var buf strings.Builder
	buf.WriteRune('[')
	for _, it := range its {
		buf.WriteString(it.String())
	}
	buf.WriteString("],")
	return buf.String()
}

func NewValueItem(v int) *item {
	return &item{isValue: true, value: v}
}

func NewListItem(l []*item) *item {
	return &item{list: l}
}

// implements heap.Interface

func (its items) Len() int { return len(its) }

func (its items) Less(i, j int) bool {
	c, err := its[i].compareItems(its[j])
	if err != nil {
		panic("Failed comparision")
	}
	return c < 0
}

func (pq items) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *items) Push(x any) {
	*pq = append(*pq, x.(*item))
}

func (pq *items) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil // avoid memory leak
	*pq = old[0 : n-1]
	return item
}

// end of heap.Interface

// compare items and lists

func (it *item) compareItems(other *item) (int, error) {
	if it.isValue && other.isValue {
		// compare values
		return it.value - other.value, nil
	}
	wrapIntInList := func(it *item) items {
		if !it.isValue {
			return it.list
		}
		return items{
			it,
		}
	}
	return wrapIntInList(it).compareItemLists(wrapIntInList(other))
}

func (its items) compareItemLists(other items) (int, error) {
	if len(its) == 0 {
		// left slice endded first
		if len(other) == 0 {
			// they are equal
			return 0, nil
		}
		return -1, nil
	}
	if len(other) == 0 {
		// right slice ended first
		return 1, nil
	}
	c, err := its[0].compareItems(other[0])
	if err != nil {
		return 1, err
	}
	if c < 0 {
		// left item less than righ
		return -1, nil
	}
	if c > 0 {
		// right item less than left
		return 1, nil
	}

	return its[1:].compareItemLists(other[1:])
}

// end compare items and lists

func parseInt(line string) (int, int, error) {
	var b strings.Builder
	for i, r := range line {
		if r >= '0' && r <= '9' {
			// numeric
			b.WriteRune(r)
			continue
		}
		// end of numeric, return what we collected so far
		v := input.MustAtoi(b.String())
		return v, i, nil
	}
	return 0, 0, errors.New("no integer value found in line " + line)
}

func parseItems(line string) (items, int, error) {
	var res items
	var idx int
	for idx < len(line) {
		switch line[idx] {
		case '[':
			// a list
			l, i, err := parseList(line[idx:])
			if err != nil {
				return nil, 0, err
			}
			idx += i
			res = append(res, l)

		case ',':
			// just skip
			idx++
			continue

		case ']':
			// end of a items
			return res, idx, nil

		default:
			// must be integer value
			v, i, err := parseInt(line[idx:])
			if err != nil {
				return nil, 0, err
			}
			idx += i
			res = append(res, NewValueItem(v))
		}
	}
	return nil, 0, errors.New("wrong list contents format " + line)
}

func parseList(line string) (*item, int, error) {
	if line[0] != '[' {
		return nil, 0, errors.New("wrong list begin " + line)
	}

	l, i, err := parseItems(line[1:])
	if err != nil {
		return nil, 0, err
	}

	if line[i+1] != ']' {
		return nil, 0, errors.New("wrong list end " + line)
	}

	return NewListItem(l), i + 2, nil
}

func readLine(scan *bufio.Scanner) (*item, bool, error) {
	for scan.Scan() {
		line := scan.Text()
		if len(line) == 0 {
			// skip empty lines
			continue
		}

		l, _, err := parseList(line)
		if err != nil {
			return nil, false, err
		}

		return l, true, nil
	}
	if err := scan.Err(); err != nil {
		return nil, false, err
	}
	// EOF
	return nil, false, nil
}

func readPair(scan *bufio.Scanner) (*item, *item, bool, error) {
	one, ok, err := readLine(scan)
	if err != nil {
		return nil, nil, false, err
	}
	if !ok {
		// EOF
		return nil, nil, false, nil
	}
	two, _, err := readLine(scan)
	if err != nil {
		return nil, nil, false, err
	}

	return one, two, true, nil
}

func orderLists(name string) (int, error) {
	file, err := os.Open(name)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	var itsHeap items
	heap.Init(&itsHeap)

	scan := bufio.NewScanner(file)
	it, ok, err := readLine(scan)
	for ok && err == nil {
		heap.Push(&itsHeap, it)
		it, ok, err = readLine(scan)
	}
	if err != nil {
		return 0, err
	}

	it2, _, _ := parseList("[[2]]")
	heap.Push(&itsHeap, it2)

	it6, _, _ := parseList("[[6]]")
	heap.Push(&itsHeap, it6)

	var idx, i2, i6 int
	for itsHeap.Len() > 0 {
		idx++
		it := heap.Pop(&itsHeap).(*item)
		if it == it2 {
			i2 = idx
			if i6 > 0 {
				// found both indexes we're looking for
				break
			}
		}
		if it == it6 {
			i6 = idx
			if i2 > 0 {
				// found both indexes we're looking for
				break
			}
		}
	}

	return i2 * i6, nil
}

func compareLists(name string) (int, error) {
	file, err := os.Open(name)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scan := bufio.NewScanner(file)
	var idx int
	var rightPairs int
	one, two, ok, err := readPair(scan)
	if err != nil {
		return 0, err
	}
	for ok {
		idx++
		b, err := one.compareItems(two)
		if err != nil {
			return 0, err
		}

		if b <= 0 {
			rightPairs += idx
		}

		one, two, ok, err = readPair(scan)
		if err != nil {
			return 0, err
		}
	}

	return rightPairs, nil
}

func main() {
	res, err := compareLists("data/part_one_short.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
	fmt.Println("=================")
	assert.Equals(13, res, "")

	res, err = compareLists("data/input.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
	fmt.Println("=================")
	assert.Equals(6395, res, "")

	res, err = orderLists("data/part_one_short.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
	fmt.Println("=================")
	assert.Equals(140, res, "")

	res, err = orderLists("data/input.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
	fmt.Println("=================")
	assert.Equals(24921, res, "")
}
