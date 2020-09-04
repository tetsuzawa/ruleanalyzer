package ruleanalyzer

type TplConfig struct {
	Name  string
	Queue MilestoneQueue
}

func NewTplConfig(name string, q MilestoneQueue) *TplConfig {
	return &TplConfig{Name: name, Queue: q}
}
