package main

import (
	"testing"
)

func TestReadFileSuccess(t *testing.T) {
	excepted := "tick"
	result := readFile()
	if result == " " {
		t.Errorf("Failed : excepted : %s but result : %s", excepted, result)
	}
}

/* Set "quack" in custome file */
func TestReadFileWithQuackInCustomFileSuccess(t *testing.T) {
	excepted := "quack"
	result := readFile()
	if result != excepted {
		t.Errorf("Failed : excepted : %s but result : %s", excepted, result)
	}
}
