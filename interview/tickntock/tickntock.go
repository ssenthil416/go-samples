package main

import (
	"fmt"
	"os"
	"time"
)

var (
	fName = "custom.txt"
	fSize = 10
	lm    = time.Now()
	tick  = "tick"
)

func main() {
	st := time.Now()
	sTicker := time.NewTicker(1 * time.Second)
	go func(st time.Time) {
		for t := range sTicker.C {
			mt := st.Add(1 * time.Minute)
			ht := st.Add(1 * time.Hour)
			if ht.Unix() == t.Unix() {
				st = ht
			} else if mt.Unix() == t.Unix() {
				st = mt
			} else {
				fmt.Println(readFile(), t)
			}
		}
	}(st)
	mTicker := time.NewTicker(1 * time.Minute)
	go func(st time.Time) {
		for t := range mTicker.C {
			ht := st.Add(1 * time.Hour)
			if ht.Unix() == t.Unix() {
				st = ht
			} else {
				fmt.Println("tock", t)
			}
		}
	}(st)
	hTicker := time.NewTicker(1 * time.Hour)
	go func() {
		for t := range hTicker.C {
			fmt.Println("bong", t)
		}
	}()
	cTicker := time.NewTicker(3 * time.Hour)
	<-cTicker.C
}

func readFile() string {
	fh, err := os.Open(fName)
	if err != nil {
		return tick
	}
	fi, err := fh.Stat()
	fs := fi.Size()
	if err != nil || fs == 0 || fs > 10 {
		return tick
	}

	b := make([]byte, fs)
	_, err = fh.Read(b)
	if err == nil {
		return string(b[:fs-1])
	}

	return tick
}
