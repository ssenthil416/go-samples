package main

import (
	"fmt"
	"strings"
	"sync"

	"github.com/go-samples/interview/wm/spellcheck/dic"
	"github.com/go-samples/interview/wm/spellcheck/input"
)

const (
	numOfWorkers = 4
)

type JobInfo struct {
	lineString string
	lineNumber int
}

var (
	dictionary     map[string][]string
	stopWorkerChan chan struct{}
	jobChan        chan JobInfo
	jobDoneChan    chan bool
)

func main() {
	// allocate dic
	dictionary = make(map[string][]string)

	err := dic.GetDicReady(dictionary)
	if err != nil {
		fmt.Printf("Error :%+v\n", err)
		return
	}

	inLines, err := input.GetLinesToValidate()
	if err != nil {
		fmt.Printf("Error :%+v\n", err)
	}

	//fmt.Println("Lines to validate :", len(inLines))

	// init Chan
	stopWorkerChan = make(chan struct{})
	jobChan = make(chan JobInfo, numOfWorkers)
	jobDoneChan = make(chan bool, len(inLines))

	// start worker
	wg := &sync.WaitGroup{}
	for i := 0; i < numOfWorkers; i++ {
		wg.Add(1)
		go runWorker(wg, jobChan, stopWorkerChan, jobDoneChan)
	}

	// assign job to worker
	for lineNum, lineStr := range inLines {
		tmp := JobInfo{}
		tmp.lineString = lineStr
		tmp.lineNumber = lineNum + 1
		jobChan <- tmp
	}

	// validate all job is done
	count := 1
	for {
		if <-jobDoneChan {
			count++
			if count == len(inLines) {
				break
			}
		}
	}

	// time to close job channel
	close(jobChan)

	// time stop all workers
	close(stopWorkerChan)

	// wait for all worker to close
	wg.Wait()
}

func runWorker(wg *sync.WaitGroup, jobChan chan JobInfo, stopWorkerChan chan struct{}, jobDoneChan chan<- bool) {
	for {
		select {
		case <-stopWorkerChan:
			wg.Done()
			return
		case ji := <-jobChan:

			wsfl := strings.Split(ji.lineString, " ")
			cn := 1
			for _, w := range wsfl {
				lw := strings.ToLower(w)

				// avoid word which got number and special char
				if avoidWord([]byte(lw)) {
					cn = cn + len(lw) + 1
					continue
				}

				//fmt.Println(" Before extraAToZOnly :", string(lw))

				// extra only atoz in a word
				vw := extraAToZOnly([]byte(lw))

				// avoid word which got number and special char
				if avoidWord(vw) {
					cn = cn + len(lw) + 1
					continue
				}

				sw, ok := validate(vw)
				if !ok {
					fmt.Printf("Line Number :%d, Column Number :%d, Wrong word :%s, Suggested Word:%s\n", ji.lineNumber, cn, vw, sw)
				}
				cn = cn + len(lw) + 1
			}
			jobDoneChan <- true
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
