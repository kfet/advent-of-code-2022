package main

import (
	"container/heap"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"kfet.org/aoc_common/assert"
	"kfet.org/aoc_common/input"
)

type entryType uint8

const (
	dir = iota
	file
)

type entry struct {
	t    entryType
	name string
	size int

	parent   *entry
	children map[string]*entry
}

func (e *entry) enumerate(level int, visit func(*entry, int) bool) {
	if !visit(e, level) {
		return
	}
	if e.t == file {
		// done
		return
	}
	for _, ch := range e.children {
		ch.enumerate(level+1, visit)
	}
}

func (e *entry) StringWithIndent(indent string) string {

	ind := func(level int) string {
		var sb strings.Builder
		sb.WriteString(indent)
		for i := 0; i < level; i++ {
			sb.WriteString("  ")
		}
		return sb.String()
	}

	var sb strings.Builder
	e.enumerate(0, func(e1 *entry, level int) bool {
		if e.t == file {
			sb.WriteString(fmt.Sprintf("%s%s %d\n", ind(level), e.name, e.size))
			return true
		}
		sb.WriteString(fmt.Sprintf("%s(dir) %s %d\n", ind(level), e.name, e.size))
		return true
	})
	return sb.String()
}

func (e *entry) String() string {
	return e.StringWithIndent("")
}

func NewEntry(parent *entry, t entryType, name string, size int) *entry {
	e := &entry{
		t:      t,
		name:   name,
		size:   size,
		parent: parent,
	}
	e.children = make(map[string]*entry)
	return e
}

type fs struct {
	root        *entry
	pwd         *entry
	dirSizeHook func(e *entry)
}

func NewFs() *fs {
	r := NewEntry(nil, dir, "/", 0)
	f := &fs{
		root: r,
		pwd:  r,
	}
	return f
}

func (f *fs) updateParentsSize(p *entry, size int) {
	parent := p
	for parent != nil {
		parent.size += size
		if f.dirSizeHook != nil {
			f.dirSizeHook(parent)
		}
		parent = parent.parent
	}
}

func (f *fs) mkfile(args []string) error {
	if len(args) != 2 {
		return errors.New("wrong mkfile args " + fmt.Sprint(args))
	}
	size, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}
	name := args[1]
	if entry, ok := f.pwd.children[name]; ok {
		// already exists
		if entry.t != file {
			return errors.New("directory already exists with name " + name)
		}
		if entry.size != size {
			return errors.New("entry exists with name " + name + " and size " + args[0])
		}
		// entry matches, noop
		return nil
	}

	entry := NewEntry(f.pwd, file, name, size)
	f.pwd.children[name] = entry

	// count the entry size in all parent directories
	f.updateParentsSize(entry.parent, entry.size)

	return nil
}

func (f *fs) mkdir(name string) error {
	if d, ok := f.pwd.children[name]; ok {
		if d.t == dir {
			// already exists
			return nil
		}
		return errors.New("cannot make dir, file name already exists " + name)
	}
	e := NewEntry(f.pwd, dir, name, 0)
	f.pwd.children[name] = e
	if f.dirSizeHook != nil {
		f.dirSizeHook(e)
	}
	return nil
}

func (f *fs) cd(name string) error {
	switch name {
	case "/":
		f.pwd = f.root
	case "..":
		f.pwd = f.pwd.parent
	default:
		entry, ok := f.pwd.children[name]
		if !ok {
			return errors.New("dir not found " + name)
		}
		if entry.t != dir {
			return errors.New("cannot cd into a file " + name)
		}
		f.pwd = entry
	}
	return nil
}

func (f *fs) execLine(line string) error {
	args := strings.Split(line, " ")
	if len(args) < 2 {
		return errors.New("too few entries in line " + line)
	}
	switch {
	case args[0] == "$" && args[1] == "cd":
		return f.cd(args[2])
	case args[0] == "$" && args[1] == "ls":
		// noop
	case args[0] == "dir":
		return f.mkdir(args[1])
	default:
		return f.mkfile(args)
	}
	return nil
}

// An IntHeap is a min-heap of ints.
type DirHeap []*entry

func (h DirHeap) Len() int           { return len(h) }
func (h DirHeap) Less(i, j int) bool { return h[i].size < h[j].size }
func (h DirHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *DirHeap) Push(x any)        { *h = append(*h, x.(*entry)) }

func (h *DirHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func buildDirHeap(f *fs) *DirHeap {
	h := &DirHeap{}
	heap.Init(h)

	f.root.enumerate(0, func(e *entry, i int) bool {
		if e.t == dir {
			heap.Push(h, e)
		}
		return true
	})
	return h
}

func findMatchingDir(f *fs) (int, error) {
	h := buildDirHeap(f)

	needSpace := 30_000_000 - (70_000_000 - f.root.size)
	for h.Len() > 0 {
		d := heap.Pop(h).(*entry)
		if d.size > needSpace {
			return d.size, nil
		}
	}

	return 0, errors.New("no direcotry is big enough")
}

func findSmallDirs(f *fs) (int, error) {
	h := buildDirHeap(f)

	var sum int
	for h.Len() > 0 {
		d := heap.Pop(h).(*entry)
		if d.size <= 100_000 {
			sum += d.size
		} else {
			break
		}
	}
	return sum, nil
}

func processFile(name string, resCalc func(*fs) (int, error)) (int, error) {
	f := NewFs()

	err := input.ReadFileLines(name, func(line string) error {
		return f.execLine(line)
	})
	if err != nil {
		return 0, err
	}

	return resCalc(f)
}

func main() {
	res, err := processFile("data/part_one.txt", findSmallDirs)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
	fmt.Println("=================")
	assert.Equals(95437, res, "")

	res, err = processFile("data/input.txt", findSmallDirs)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
	fmt.Println("=================")
	assert.Equals(1845346, res, "")

	res, err = processFile("data/part_one.txt", findMatchingDir)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
	fmt.Println("=================")
	assert.Equals(24933642, res, "")

	res, err = processFile("data/input.txt", findMatchingDir)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
	fmt.Println("=================")
	assert.Equals(3636703, res, "")
}
