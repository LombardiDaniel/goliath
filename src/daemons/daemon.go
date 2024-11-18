package daemons

import (
	"log/slog"
	"time"
)

type Task struct {
	Interval time.Duration
	Callable func() error
}

type TaskRunner struct {
	tasks []Task
}

func (f *TaskRunner) RegisterTask(interval time.Duration, callable func() error) {
	f.tasks = append(f.tasks, Task{
		Interval: interval,
		Callable: callable,
	})
}

func (f *TaskRunner) Run() {
	for _, v := range f.tasks {
		go func(t Task) {
			for {
				err := v.Callable()
				if err != nil {
					slog.Error(err.Error())
				}
				time.Sleep(t.Interval)
			}
		}(v)
	}
}
