package ruleanalyzer

import "fmt"

func generate(q MilestoneQueue) error {
	for i, v := range q {
		fmt.Println(i, v)
	}
	return nil
}
