package dagger

import (
	"github.com/facebookgo/errgroup"
)

// We intentionally don't expose the list of dependencies on a task, with the idea being that if you
// can only set dependencies at the time of construction it's harder (impossible?) to cause
// dependency cycles.

// A Task represents a function with some dependent tasks.
type Task struct {
	deps   []*Task
	action func() error
	done   chan struct{}
	err    error
}

// NewTask creates a task from the given action function and dependencies.
func NewTask(action func() error, deps ...*Task) *Task {
	return &Task{action: action, deps: deps}
}

// Execute takes in a list of Tasks and executes each one in its own goroutine, as soon as its
// dependent Tasks have finished. It is not guaranteed that all Tasks will execute. If a Task
// returns an error, Tasks that depend on it will not execute.
// Execute will never return if it is impossible for any task to
// complete (either because there are cycles in your dependency graph or because a dependent task
// wasn't passed to Execute).
func Execute(tasks ...*Task) error {
	for _, task := range tasks {
		// Give every task a done channel that will be closed when the task completes.
		task.done = make(chan struct{})
	}
	wg := errgroup.Group{}
	for _, task := range tasks {
		wg.Add(1)
		go func(task *Task) {
			defer close(task.done)
			defer wg.Done()
			if err := executeSingleTask(task, task.deps...); err != nil {
				task.err = err
				wg.Error(err)
			}
		}(task)
	}
	return wg.Wait()
}

func executeSingleTask(task *Task, deps ...*Task) error {
	if err := waitForTasks(deps...); err != nil {
		// If a dependency errored, set an error on this task as well so that further dependencies
		// don't try to run.
		task.err = err
		return nil
	}
	return task.action()
}

func waitForTasks(tasks ...*Task) error {
	for _, task := range tasks {
		<-task.done
		if task.err != nil {
			return task.err
		}
	}
	return nil
}
