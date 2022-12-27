package main

import (
	"container/heap"
	"errors"
	"fmt"
	"math"
	"strings"

	"github.com/samber/lo"
	"kfet.org/aoc_common/assert"
	"kfet.org/aoc_common/calc"
	"kfet.org/aoc_common/input"
)

type blizzardRing struct { // used for 'left' and 'up' rings
	len int
	m   map[int]struct{}
}

func (br *blizzardRing) hasBlizzard(i, time int) bool {
	i += time
	i %= br.len
	return calc.IsPresent(br.m, i)
}

type rightBlizzardRing blizzardRing // used for 'righ' and 'down' rings

func (rbr *rightBlizzardRing) hasBlizzard(i, time int) bool {
	i -= time
	i %= rbr.len
	i += rbr.len
	i %= rbr.len
	return calc.IsPresent(rbr.m, i)
}

type field struct {
	w, h int                  // field width/height
	lcd  int                  // width and height lowest common denominator
	s, e int                  // field start/exit X-index
	lbr  []*blizzardRing      // '<' rings
	rbr  []*rightBlizzardRing // '>' rings
	ubr  []*blizzardRing      // '^' rings
	dbr  []*rightBlizzardRing // 'v' rings
}

func (f *field) String(exp pos, dest goal) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintln("w, h, lcd: ", f.w, f.h, f.lcd))
	sb.WriteString(fmt.Sprintln("s, e: ", f.s, f.e))
	sb.WriteString(fmt.Sprintln("pos: ", exp))
	sb.WriteString(fmt.Sprintln("dist: ", exp.calcDist(dest)))

	for y := 0; y < f.h; y++ {
		for x := 0; x < f.w; x++ {

			if x == exp.x && y == exp.y {
				sb.WriteRune('E')
				continue
			}

			var r rune
			var n int
			if f.lbr[y].hasBlizzard(x, exp.t) {
				r = '<'
				n++
			}
			if f.rbr[y].hasBlizzard(x, exp.t) {
				r = '>'
				n++
			}
			if f.ubr[x].hasBlizzard(y, exp.t) {
				r = '^'
				n++
			}
			if f.dbr[x].hasBlizzard(y, exp.t) {
				r = 'v'
				n++
			}

			switch n {
			case 0:
				sb.WriteRune('.')
			case 1:
				sb.WriteRune(r)
			default:
				sb.WriteString(fmt.Sprint(n))
			}
		}
		sb.WriteRune('\n')
	}

	return sb.String()
}

func NewField(matrix []string) *field {
	f := &field{}

	f.w = len(matrix[0]) - 2
	f.h = len(matrix) - 2
	f.lcd = calc.LCD(f.w, f.h)

	f.s = strings.IndexRune(matrix[0], '.') - 1
	f.e = strings.IndexRune(matrix[f.h+1], '.') - 1

	// Create all rings
	f.lbr = make([]*blizzardRing, 0)
	f.rbr = make([]*rightBlizzardRing, 0)
	f.ubr = make([]*blizzardRing, 0)
	f.dbr = make([]*rightBlizzardRing, 0)

	for x := 0; x < f.w; x++ {
		f.ubr = append(f.ubr, &blizzardRing{m: map[int]struct{}{}, len: f.h})
		f.dbr = append(f.dbr, &rightBlizzardRing{m: map[int]struct{}{}, len: f.h}) // 'v' ring same as '>' ring
	}
	for y := 0; y < f.h; y++ {
		f.lbr = append(f.lbr, &blizzardRing{m: map[int]struct{}{}, len: f.w})
		f.rbr = append(f.rbr, &rightBlizzardRing{m: map[int]struct{}{}, len: f.w})
	}

	// Read all rings from the matrix
	f.readMatrix(matrix)

	return f
}

func (f *field) readMatrix(matrix []string) error {

	for y := 1; y < len(matrix)-1; y++ {
		for x := 1; x < len(matrix[0])-1; x++ {
			r := matrix[y][x]
			switch r {
			case '<':
				f.lbr[y-1].m[x-1] = struct{}{}
			case '>':
				f.rbr[y-1].m[x-1] = struct{}{}
			case '^':
				f.ubr[x-1].m[y-1] = struct{}{}
			case 'v':
				f.dbr[x-1].m[y-1] = struct{}{}
			case '.':
			default:
				return errors.New(fmt.Sprint("Wonrg rune ", r, " in line ", matrix[y]))
			}
		}
	}

	return nil
}

func (f *field) hasBlizzard(p pos) bool {
	return f.lbr[p.y].hasBlizzard(p.x, p.t) ||
		f.rbr[p.y].hasBlizzard(p.x, p.t) ||
		f.ubr[p.x].hasBlizzard(p.y, p.t) ||
		f.dbr[p.x].hasBlizzard(p.y, p.t)
}

type fieldState struct {
	exp       pos // expedition position - x, y, t (time)
	dist      int // distance to goal - calculated after position is determined
	timeIndex int // type modulus the world loop length, i.e. field state repeat
	f         *field
}

func (fs *fieldState) String(dest goal) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintln("time index: ", fs.timeIndex))
	sb.WriteString(fs.f.String(fs.exp, dest))
	return sb.String()
}

func (p pos) calcDist(dest goal) int {
	return calc.TaxiCab(p.x, p.y, dest.x, dest.y) + p.t
}

func NewFieldState(f *field, dest goal) *fieldState {
	exp := pos{
		x: f.s,
		y: -1,
		t: 0,
	}
	return &fieldState{
		exp:  exp,
		dist: exp.calcDist(dest),
		f:    f,
	}
}

func (fs *fieldState) copy(dp pos, dest goal) *fieldState {
	exp := pos{
		x: fs.exp.x + dp.x,
		y: fs.exp.y + dp.y,
		t: fs.exp.t + dp.t,
	}

	return &fieldState{
		exp:       exp,
		dist:      exp.calcDist(dest),
		timeIndex: exp.t % fs.f.lcd,
		f:         fs.f,
	}
}

type timeSpace struct {
	x, y, timeIndex int
}

// Here we track visited positions to deted loop
var visited = map[timeSpace]struct{}{}

func (fs *fieldState) allowed() bool {

	ts := timeSpace{fs.exp.x, fs.exp.y, fs.timeIndex}
	if _, ok := visited[ts]; ok {
		// loop detected
		return false
	}
	visited[ts] = struct{}{}

	// Allow initial states
	if fs.exp.x == fs.f.s && fs.exp.y == -1 {
		return true
	}

	// Allow end state
	if fs.exp.x == fs.f.e && fs.exp.y == fs.f.h {
		return true
	}

	if fs.exp.x < 0 || fs.exp.x > fs.f.w-1 ||
		fs.exp.y < 0 || fs.exp.y > fs.f.h-1 {
		// outside of bounds
		return false
	}

	if fs.f.hasBlizzard(fs.exp) {
		// there's a blizzard on this spot
		return false
	}

	return true
}

type pos struct {
	x, y, t int
}

func (p pos) String() string {
	return fmt.Sprint(p.x, p.y, p.t)
}

func (fs *fieldState) nextStates(dest goal) []*fieldState {
	return lo.FilterMap([]pos{
		{0, 0, 1},  // stay put
		{-1, 0, 1}, // left
		{+1, 0, 1}, // right
		{0, -1, 1}, // up
		{0, +1, 1}, // down
	}, func(item pos, index int) (*fieldState, bool) {
		res := fs.copy(item, dest)
		return res, res.allowed()
	})
}

type fsHeap []*fieldState

func (fh fsHeap) Len() int           { return len(fh) }
func (fh fsHeap) Less(i, j int) bool { return fh[i].dist < fh[j].dist }
func (fh fsHeap) Swap(i, j int)      { fh[i], fh[j] = fh[j], fh[i] }
func (fh *fsHeap) Push(x any)        { *fh = append(*fh, x.(*fieldState)) }
func (fh *fsHeap) Pop() any {
	old := *fh
	n := len(old)
	x := old[n-1]
	*fh = old[0 : n-1]
	return x
}

func (fh *fsHeap) pushAll(fss []*fieldState) {
	for _, fs := range fss {
		heap.Push(fh, fs)
	}
}

func (fs *fieldState) minPathGoals(maxTime int, goals []goal) (*fieldState, bool) {
	nfs := fs
	for _, g := range goals {
		// reset visited map for each new goal
		visited = map[timeSpace]struct{}{}
		nfs.dist = nfs.exp.calcDist(g)

		var found bool
		nfs, found = nfs.minPathBFS(maxTime, g)
		if !found {
			return nil, false
		}
	}

	return nfs, true
}

func (fs *fieldState) minPathBFS(maxTime int, dest goal) (*fieldState, bool) {

	fh := &fsHeap{fs}
	heap.Init(fh)

	nextFh := &fsHeap{}
	heap.Init(nextFh)

	for len(*fh) > 0 {
		for len(*fh) > 0 {
			ns := heap.Pop(fh).(*fieldState)

			if ns.exp.x == dest.x && ns.exp.y == dest.y {
				return ns, true
			}

			if ns.dist < maxTime {
				nextFh.pushAll(ns.nextStates(dest))
			}
		}

		fh = nextFh

		nextFh = &fsHeap{}
		heap.Init(nextFh)
	}

	// No path found
	return nil, false
}

type goal struct {
	x, y int
}

func processFile(fileName string, goalNum int) (int, error) {

	fmt.Println("Processing ", fileName)

	matrix := []string{}
	err := input.ReadFileLines(fileName, func(line string) error {
		matrix = append(matrix, line)
		if len(matrix) > 1 {
			if len(line) != len(matrix[len(matrix)-2]) {
				return errors.New(fmt.Sprint("Wrong line lenggh ", line))
			}
		}
		return nil
	})
	if err != nil {
		return 0, err
	}

	f := NewField(matrix)
	if f.s < 0 || f.e < 0 {
		return 0, errors.New(fmt.Sprint("Start or end not found", f.s, f.e))
	}

	goals := []goal{
		{f.e, f.h},
		{f.s, -1},
		{f.e, f.h},
	}
	if goalNum > len(goals) {
		return 0, errors.New(fmt.Sprint("invalid number of goals ", goalNum))
	}

	fs := NewFieldState(f, goals[0])
	res, found := fs.minPathGoals(math.MaxInt, goals[0:goalNum])
	if !found {
		return 0, errors.New("path not found")
	}

	return res.exp.t, nil
}

func main() {
	res, err := processFile("data/part_one.txt", 1)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
	fmt.Println("=================")
	assert.Equals(10, res, "")

	res, err = processFile("data/part_one_two.txt", 1)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
	fmt.Println("=================")
	assert.Equals(18, res, "")

	res, err = processFile("data/input.txt", 1)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
	fmt.Println("=================")
	assert.Equals(290, res, "")

	res, err = processFile("data/part_one.txt", 3)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
	fmt.Println("=================")
	assert.Equals(30, res, "")

	res, err = processFile("data/part_one_two.txt", 3)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
	fmt.Println("=================")
	assert.Equals(54, res, "")

	res, err = processFile("data/input.txt", 3)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
	fmt.Println("=================")
	assert.Equals(842, res, "")
}
