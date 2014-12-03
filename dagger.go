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
}

// NewTask creates a task from the given action function and dependencies.
func NewTask(action func() error, deps ...*Task) *Task {
	return &Task{action: action, deps: deps}
}

// Execute takes in a list of tasks and executes them as soon as possible. Execute will never return
// if it is impossible for any task to complete (either because there are cycles in your dependency
// graph or because a dependent task wasn't passed to Execute).
func Execute(tasks ...*Task) error {
	for _, task := range tasks {
		// Give every task a done channel
		task.done = make(chan struct{})
	}
	wg := errgroup.Group{}
	for _, task := range tasks {
		wg.Add(1)
		go func(task *Task) {
			defer close(task.done)
			defer wg.Done()
			// Wait for all dependencies to finish
			for _, dep := range task.deps {
				<-dep.done
			}
			if err := task.action(); err != nil {
				wg.Error(err)
			}
		}(task)
	}
	return wg.Wait()
}
