package seed

// TaskAble ...
type TaskAble interface {
	Call(*Process, *Task) error
}

// Task ...
type Task struct {
	TaskAble
	*Thread
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
func NewTask(task TaskAble) *Task {
	tsk := new(Task)
	tsk.Thread = NewThread()
	tsk.TaskAble = task
	return tsk
}
