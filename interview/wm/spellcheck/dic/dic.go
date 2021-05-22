package dic

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func GetDicReady(dictionary map[string][]string) error {

	// opens specific file in read-only
	file, err := os.Open("./dic/dictionary.txt")
	if err != nil {
		return fmt.Errorf("dic file open failed :%v", err)
	}
	defer file.Close()

	// read content from the file
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var lines []string

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	//fmt.Printf("num of dictionary lines :%d\n", len(lines))

	// fill inmemory dic for validation
	for _, w := range lines {
		lw := strings.ToLower(w)
		k := string([]byte(lw[:1]))
		if val, ok := dictionary[k]; !ok {
			val := append(val, lw)
			dictionary[k] = val
		} else {
			val = append(val, lw)
			dictionary[k] = val
		}
	}
	//fmt.Printf("len of dictionary  :%d\n", len(dictionary))
	return nil
}
