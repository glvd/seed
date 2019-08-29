package seed

import "context"

// Task ...
type Task struct {
	*Thread
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
