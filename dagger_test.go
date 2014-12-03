package dagger

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTrivial(t *testing.T) {
	task := NewTask(func() error { return nil })
	assert.Nil(t, Execute(task))
}

func TestTrivialError(t *testing.T) {
	task := NewTask(func() error { return errors.New("an error") })
	assert.EqualError(t, Execute(task), "an error")
}

func TestSeries(t *testing.T) {
	// Please don't ever actually use dagger to synchronize shared memory like this.
	// Use channels for communication between goroutines, not shared memory.
	i := 0
	t1 := NewTask(func() error { return nil })
	t2 := NewTask(func() error {
		i = 1
		return nil
	}, t1)
	t3 := NewTask(func() error {
		assert.Equal(t, i, 1)
		return nil
	}, t2)
	assert.Nil(t, Execute(t1, t2, t3))
}

func TestSeriesError(t *testing.T) {
	// Please don't ever actually use dagger to synchronize shared memory like this.
	// Use channels for communication between goroutines, not shared memory.
	i := 0
	t1 := NewTask(func() error { return errors.New("an error") })
	t2 := NewTask(func() error {
		i = 1
		return nil
	}, t1)
	assert.EqualError(t, Execute(t1, t2), "an error")
	assert.Equal(t, i, 0)
}

func TestMultiError(t *testing.T) {
	t1 := NewTask(func() error { return errors.New("error1") })
	t2 := NewTask(func() error { return errors.New("error2") })
	assert.EqualError(t, Execute(t1, t2), "multiple errors: error1 | error2")
}

func ExampleExecute() {
	doNothing := func() error { return nil }
	t1 := NewTask(doNothing)
	t2 := NewTask(doNothing, t1)
	t3 := NewTask(doNothing, t1)
	t4 := NewTask(doNothing, t1, t2)
	t5 := NewTask(doNothing, t4)
	t6 := NewTask(doNothing, t1)
	t7 := NewTask(doNothing, t5, t6)
	Execute(t1, t2, t3, t4, t5, t6, t7)
	// Output:
}
