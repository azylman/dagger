# dagger
--
    import "gopkg.in/azylman/dagger.v1"

Package dagger provides a way to define a set of tasks with dependent tasks and
a way to execute those tasks in the most optimal way - that is, begin executing
each task as soon as its dependent tasks have completed.

## Usage

#### func  Execute

```go
func Execute(tasks ...*Task) error
```
Execute takes in a list of Tasks and executes each one in its own goroutine, as
soon as its dependent Tasks have finished. It is not guaranteed that all Tasks
will execute. If a Task returns an error, Tasks that depend on it will not
execute. Execute will never return if it is impossible for any task to complete
(either because there are cycles in your dependency graph or because a dependent
task wasn't passed to Execute).

#### type Task

```go
type Task struct {
}
```

A Task represents a function with some dependent tasks.

#### func  NewTask

```go
func NewTask(action func() error, deps ...*Task) *Task
```
NewTask creates a task from the given action function and dependencies.
