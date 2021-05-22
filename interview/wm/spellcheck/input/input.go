package input

import "github.com/go-samples/interview/wm/spellcheck/common"

func GetLinesToValidate() ([]string, error) {
	return common.ReadLinesFromFile("./input/big.txt")
}
