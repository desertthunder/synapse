package main

type Settings struct {
	MaxRetries uint32
}

type Scheduler struct {
	Queue    []Task
	Settings Settings
}

type Task interface {
	Action(data any) error
}

// CronJob represents a recurring task
type CronJob struct{}

// Operation represents a single instance of a task
type Operation struct{}

// function TaskRunner dispatches goroutines to execute tasks
func TaskRunner() error {
	return nil
}
