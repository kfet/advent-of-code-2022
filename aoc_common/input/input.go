package input

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	"github.com/samber/lo"
)

func CopyMap[K comparable, V any](source map[K]V, filter func(key K, value V) bool) map[K]V {
	res := make(map[K]V)
	for k, v := range source {
		if filter(k, v) {
			res[k] = v
		}
	}
	return res
}

func MustAtoi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic("wrong int format " + s)
	}
	return i
}

func MustAtoInts(stringInts []string) []int {
	return lo.Map(stringInts, func(item string, _ int) int {
		return MustAtoi(item)
	})
}

func ReadFileLinesStrings(name string, useLineStrings func(tokens []string) error) error {
	err := ReadFileLines(name, func(line string) error {
		strings := strings.Split(line, " ")
		err := useLineStrings(strings)
		return err
	})
	return err
}

func ReadFileLines(name string, useLine func(line string) error) error {
	file, err := os.Open(name)
	if err != nil {
		return err
	}
	defer file.Close()

	scan := bufio.NewScanner(file)
	for scan.Scan() {
		txt := scan.Text()
		err = useLine(txt)
		if err != nil {
			return err
		}
	}

	if scan.Err() != nil {
		return scan.Err()
	}

	return nil
}
