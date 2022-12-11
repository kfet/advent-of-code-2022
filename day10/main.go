package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"kfet.org/aoc_common/assert"
	"kfet.org/aoc_common/input"
)

type vmState struct {
	reg_x     int64
	cycles    int64
	traceHook func(*vmState)
}

func (vm *vmState) trace() {
	if vm.traceHook != nil {
		vm.traceHook(vm)
	}
}

type cpuOp func(ctx *vmState, args []string) error

var instructionSet = map[string]cpuOp{
	"noop": func(vm *vmState, args []string) error {
		vm.cycles++
		vm.trace()
		return nil
	},
	"addx": func(vm *vmState, args []string) error {
		if len(args) != 1 {
			return errors.New("Malformed instruction addx " + strings.Join(args, " "))
		}
		v, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return err
		}

		vm.cycles++
		vm.trace()

		vm.cycles++
		vm.trace()

		vm.reg_x += v

		return nil
	},
}

func NewVm(traceHook func(*vmState)) *vmState {
	return &vmState{
		reg_x:     1,
		traceHook: traceHook,
	}
}

func (vm *vmState) exec(file string) error {
	err := input.ReadFileLinesStrings(file, func(tokens []string) error {
		if len(tokens) == 0 {
			// skip empty line
			return nil
		}

		err := instructionSet[tokens[0]](vm, tokens[1:])
		if err != nil {
			return err
		}

		return nil
	})
	return err
}

func (vm *vmState) signalStrength() int64 {
	return vm.reg_x * vm.cycles
}

func main() {
	vm := NewVm(func(vm *vmState) {
		fmt.Println(vm.signalStrength())
	})
	err := vm.exec("data/part_one_short.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("=============")

	var signalTotal int64
	signalStrengthTracer := func(vm *vmState) {
		if vm.cycles < 20 {
			return
		}

		if (vm.cycles-20)%40 == 0 {
			signalTotal += vm.signalStrength()
			fmt.Println(vm.signalStrength())
		}
	}

	vm = NewVm(signalStrengthTracer)
	err = vm.exec("data/part_one.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(signalTotal)
	fmt.Println("=============")
	assert.Equals(int64(13140), signalTotal, "")

	signalTotal = 0
	vm = NewVm(signalStrengthTracer)
	err = vm.exec("data/input.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(signalTotal)
	fmt.Println("=============")
	assert.Equals(int64(12460), signalTotal, "")

	var screenBuilder strings.Builder
	crtTrancer := func(vm *vmState) {
		posX := (vm.cycles - 1) % 40
		if posX == 0 {
			screenBuilder.WriteString("\n")
		}

		if posX >= vm.reg_x-1 && posX <= vm.reg_x+1 {
			screenBuilder.WriteString("#")
		} else {
			screenBuilder.WriteString(".")
		}
	}

	vm = NewVm(crtTrancer)
	err = vm.exec("data/part_one.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(screenBuilder.String())
	fmt.Println("=============")

	expected := `
##..##..##..##..##..##..##..##..##..##..
###...###...###...###...###...###...###.
####....####....####....####....####....
#####.....#####.....#####.....#####.....
######......######......######......####
#######.......#######.......#######.....`
	assert.Equals(expected, screenBuilder.String(), "")

	screenBuilder.Reset()
	vm = NewVm(crtTrancer)
	err = vm.exec("data/input.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(screenBuilder.String())
	fmt.Println("=============")

	expected = `
####.####.####.###..###...##..#..#.#....
#.......#.#....#..#.#..#.#..#.#.#..#....
###....#..###..#..#.#..#.#..#.##...#....
#.....#...#....###..###..####.#.#..#....
#....#....#....#....#.#..#..#.#.#..#....
####.####.#....#....#..#.#..#.#..#.####.`
	assert.Equals(expected, screenBuilder.String(), "")
}
