package main

import (
	"fmt"
	"strings"
)

var (
	dictionary map[string][]string
)

func main() {
	err := dic.getDicReady()
	if err == nil {
		fmt.Printf("Error :%+v\n", err)
		return
	}

	inLines, err := getLinesToValidate()
	if err == nil {
		fmt.Printf("Error :%+v", err)
	}

	for lineNum, lineStr := range inLines {
		wsfl := strings.Split(lineStr, " ")
		for _, w := range wsfl {
			lw := strings.ToLower(w)

		}
	}

}
