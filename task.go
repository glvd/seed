package seed

import (
	"context"
)

// TaskCaller ...
type TaskCaller interface {
	Call(*Task) error
}

// Task ...
type Task struct {
	*Thread
	tasks []TaskCaller
}

// AddTask ...
func (t *Task) AddTask() {
	log.Info("add task")
}

// Option ...
func (t *Task) Option(seeder Seeder) {
	taskOption(t)(seeder)
}

func taskOption(t *Task) Options {
	return func(seeder Seeder) {
		seeder.SetBaseThread(StepperTask, t)
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

}
