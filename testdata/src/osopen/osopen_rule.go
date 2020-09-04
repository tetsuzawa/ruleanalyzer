package osopen

import (
	"os"
)

func RuleOsOpen() {
	// step: call os.Open
	f, _ := os.Open("xxx")
	// step: call *File.close
	defer f.Close()
}
