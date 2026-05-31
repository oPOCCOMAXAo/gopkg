package tasks

import "time"

type registeredTask struct {
	Name       string
	Task       Task
	Timeout    time.Duration
	Dependents []string
}
