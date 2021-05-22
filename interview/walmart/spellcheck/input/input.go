package main

import (
	"bufio"
	"fmt"
	"os"
)

func getLinesToValidate() ([]string, error) {
	// opens specific file in read-only
	file, err := os.Open("dictionary.txt")
	if err != nil {
		return nil, fmt.Errorf("failed to open : dictionary.txt")
	}
	defer file.Close()

	// read content from the file
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var lines []string

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, nil
}
