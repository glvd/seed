package seed

// TaskStep ...
type TaskStep int

// TaskStep ...
const (
	TaskNone TaskStep = iota
	TaskInformation
	TaskMax
)

// TaskAble ...
type TaskAble interface {
	CallTask(*Task, *Process) error
}

// Task ...
type Task struct {
	TaskAble
	*Thread
	Step TaskStep
}

// Call ...
func (t *Task) Call(process *Process) error {
	return t.TaskAble.CallTask(t, process)
}

// NewTask ...
func NewTask(task TaskAble) *Task {
	tsk := new(Task)
	tsk.Thread = NewThread()
	tsk.TaskAble = task
	return tsk
}
