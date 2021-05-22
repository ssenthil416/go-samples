package main

import (
	"fmt"
	"strings"

	"github.com/go-samples/interview/walmart/spellcheck/dic"
	"github.com/go-samples/interview/walmart/spellcheck/input"
)

var (
	dictionary map[string][]string
)

func main() {
	err := dic.GetDicReady(dictionary)
	if err == nil {
		fmt.Printf("Error :%+v\n", err)
		return
	}

	inLines, err := input.GetLinesToValidate()
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
