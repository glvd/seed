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
	ct TaskAble
	*Thread
	Step TaskStep
}

// Call ...
func (t *Task) Call(process *Process) error {
	return t.ct.CallTask(t, process)
}

// NewTask ...
func NewTask(task TaskAble) *Task {
	tsk := new(Task)
	tsk.Thread = NewThread()
	tsk.ct = task
	return tsk
}
