package dic

import (
	"strings"

	"github.com/go-samples/interview/wm/spellcheck/common"
)

func GetDicReady(dictionary map[string][]string) error {

	lines, err := common.ReadLinesFromFile("./dic/dictionary.txt")
	if err != nil {
		return err
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
