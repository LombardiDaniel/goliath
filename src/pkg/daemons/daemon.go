package daemons

import (
	"fmt"
	"log/slog"
	"time"
)

// Task is a struct that defines a task to be run in the background.
type Task struct {
	Interval time.Duration
	Callable func() error
	Workers  uint32
}

// TaskRunner runs Tasks specified by the user, use TaskRunner.RegisterTask() to
// specify daemons.
type TaskRunner struct {
	// List of Tasks to be run in the background
	tasks []Task
}

// RegisterTask Creates and inserts a task in the TaskRunner, runs only after Dispatch() call
func (f *TaskRunner) RegisterTask(interval time.Duration, callable func() error, workers uint32) {
	f.tasks = append(f.tasks, Task{
		Interval: interval,
		Callable: callable,
		Workers:  workers,
	})
}

// Dispatches Tasks to be run
func (f *TaskRunner) Dispatch() {
	for _, v := range f.tasks {
		for range v.Workers {
			go taskRunner(v)
		}
	}
}

func taskWrapper(t Task) {
	defer func() {
		if r := recover(); r != nil {
			slog.Error(fmt.Sprintf("Task crashed: %v, waiting 5s to restart...", r))
			time.Sleep(5 * time.Second)
		}
	}()
	err := t.Callable()
	if err != nil {
		slog.Error(err.Error())
	}
}

func taskRunner(t Task) {
	for {
		taskWrapper(t)
		time.Sleep(t.Interval)
	}
}
