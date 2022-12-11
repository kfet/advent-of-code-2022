package input

import (
	"bufio"
	"os"
	"strings"
)

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
