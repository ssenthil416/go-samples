package nouns

import (
	"os"
	"time"

	"github.com/go-samples/interview/wm/spellcheck/common"
)

const (
	fn = "./nouns/nouns.txt"
)

var (
	places      = []string{}
	lastModTime time.Time
)

func Populate() (err error) {
	places, err = common.ReadLinesFromFile(fn)
	if err != nil {
		return err
	}

	// get file status
	fs, err := os.Stat(fn)
	if err != nil {
		return err
	}
	lastModTime = fs.ModTime()
	return nil
}

func Avoid(w string) bool {
	fs, err := os.Stat(fn)
	if err != nil {
		return false
	}

	if lastModTime.After(fs.ModTime()) {
		if err := Populate(); err != nil {
			return false
		}
		lastModTime = fs.ModTime()
	}

	for _, p := range places {
		if w == p {
			return true
		}
	}
	return false
}
