package main

import (
	"fmt"
	"github.com/tetsuzawa/ruleanalyzer"
	"os"
)

func main() {
	if err := ruleanalyzer.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to run ruleanalyzer: %v\n", err)
		os.Exit(1)
	}
}
