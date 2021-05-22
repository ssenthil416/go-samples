package main

import (
	"fmt"
	"strings"

	"github.com/go-samples/interview/wm/spellcheck/dic"
	"github.com/go-samples/interview/wm/spellcheck/input"
)

var (
	dictionary map[string][]string
)

func main() {
	// allocate dic
	dictionary = make(map[string][]string)

	err := dic.GetDicReady(dictionary)
	if err != nil {
		fmt.Printf("Error :%+v\n", err)
		return
	}

	/*
		fmt.Printf("len of dictionary  :%d\n", len(dictionary))

			for k, v := range dictionary {
				fmt.Printf(" key : %+v", k)
				fmt.Printf(" len of v :%d\n", len(v))
			}
	*/
	inLines, err := input.GetLinesToValidate()
	if err != nil {
		fmt.Printf("Error :%+v\n", err)
	}

	//fmt.Println("Lines to validate :", len(inLines))

	for lineNum, lineStr := range inLines {
		wsfl := strings.Split(lineStr, " ")
		for cn, w := range wsfl {
			lw := strings.ToLower(w)

			// avoid word which got number and special char
			if avoidWord([]byte(lw)) {
				continue
			}

			//fmt.Println(" Before extraAToZOnly :", string(lw))

			// extra only atoz in a word
			vw := extraAToZOnly([]byte(lw))

			// avoid word which got number and special char
			if avoidWord(vw) {
				continue
			}

			sw, ok := validate(vw)
			if !ok {
				fmt.Printf("Line Number :%d, Column Number :%d, Wrong word :%s, Suggested Word:%s\n", lineNum+1, cn, vw, sw)
			}
		}
	}

}

func validate(vw []byte) (string, bool) {
	//fmt.Println("validate :", string(vw))

	aw := dictionary[string(vw[0])]

	//fmt.Printf("len words in the category %s:%d\n", string(vw[0]), len(sw))
	sw := ""

	for _, w := range aw {
		if w == string(vw) {
			return sw, true
		}
	}

	return sw, false
}

func extraAToZOnly(val []byte) []byte {
	var nc []byte
	for _, c := range val {
		if 'a' <= c && c <= 'z' {
			nc = append(nc, c)
		}
	}
	return nc
}

func avoidWord(val []byte) bool {
	// if word len is one, continue
	if len(val) == 1 || len(val) == 0 {
		return true
	}

	switch d := string(val); d {
	case "i", "ii", "iii", "iv", "v", "vi", "vii", "ix", "x", "xii", "***":
		return true
	}

	return gotNumbers(val) || checkSpecialChar(val)
}

func checkSpecialChar(val []byte) bool {
	for _, c := range val {
		if '-' == c || '\'' == c || '/' == c || '[' == c || ']' == c || '.' == c {
			return true
		}
	}
	return false
}

func gotNumbers(val []byte) bool {
	for _, c := range val {
		if '0' <= c && c <= '9' {
			return true
		}
	}
	return false
}
