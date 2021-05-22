package main

import (
	"fmt"

	ms "github.com/ssenthil416/mystr"
)

func main() {
	//strMirror := "abcddcba"
	strMirror := ""
	if ms.StrMirror(strMirror) {
		fmt.Println("Yes, String is Mirror")
	} else {
		fmt.Println("No, String is Not Mirror")
	}
}
