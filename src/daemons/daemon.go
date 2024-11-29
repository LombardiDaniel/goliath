package daemons

import (
	"fmt"
	"log/slog"
	"time"
)

type Task struct {
	Interval time.Duration
	Callable func() error
	Workers  uint32
}

type TaskRunner struct {
	tasks []Task
}

func (f *TaskRunner) RegisterTask(interval time.Duration, callable func() error, workers uint32) {
	f.tasks = append(f.tasks, Task{
		Interval: interval,
		Callable: callable,
		Workers:  workers,
	})
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

func (f *TaskRunner) Dispatch() {
	for _, v := range f.tasks {
		for range v.Workers {
			go taskRunner(v)
		}
	}
}
