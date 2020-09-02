package main

import (
	"github.com/tetsuzawa/ruleanalyzer"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() { unitchecker.Main(ruleanalyzer.Analyzer) }

