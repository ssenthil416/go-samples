package common

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func ReadLinesFromFile(fn string) ([]string, error) {
	// opens specific file in read-only
	file, err := os.Open(fn)
	if err != nil {
		return nil, fmt.Errorf("input file open failed :%v\n", err)
	}
	defer file.Close()

	// read content from the file
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var lines []string
	for scanner.Scan() {
		lines = append(lines, strings.ToLower(scanner.Text()))
	}
	return lines, nil
}
