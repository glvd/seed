package seed

import (
	"context"
	"sync"
)

// TaskCaller ...
type TaskCaller interface {
	Call(*Task) error
}

// Task ...
type Task struct {
	*Thread
	taskMutex sync.RWMutex
	tasks     []TaskCaller
}

// AddTask ...
func (t *Task) AddTask(caller TaskCaller) {
	t.taskMutex.Lock()
	t.tasks = append(t.tasks, caller)
	t.taskMutex.Unlock()
}

// Option ...
func (t *Task) Option(seeder Seeder) {
	taskOption(t)(seeder)
}

func taskOption(t *Task) Options {
	return func(seeder Seeder) {
		seeder.SetThread(StepperTask, t)
	}
}

// NewTask ...
func NewTask() Threader {
	tsk := new(Task)
	tsk.Thread = NewThread()
	return tsk
}

// Run ...
func (t *Task) Run(ctx context.Context) {
	t.taskMutex.RLock()
	for i, tsk := range t.tasks {
		e := tsk.Call(t)
		if e != nil {
			log.With("index", i).Error(e)
		}
	}
	t.taskMutex.RUnlock()

}
