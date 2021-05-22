package main

import (
	"errors"
	"fmt"
)

func main() {
	in := []interface{}{1, 2, []int{3}, 4, []interface{}{5, 6, []int{7}, 8}}
	out, err := goFlat(in)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(out)
	}
}

func goFlat(si interface{}) ([]int, error) {
	out := make([]int, 0, 1)
	switch v := si.(type) {
	case []int:
		out = append(out, v...)
	case int:
		out = append(out, v)
	case []interface{}:
		for i := range v {
			o, err := goFlat(v[i])
			if err != nil {
				return nil, errors.New("Error : Not int or []int")
			}
			out = append(out, o...)
		}
	default:
		return nil, errors.New("Error : Not int")
	}
	return out, nil
}
