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
	Step() Stepper
	CallTask(*Task, Seeder) error
}

// Task ...
type Task struct {
	ct       TaskAble
	taskStep TaskStep
	//*Thread
}

// Push ...
func (t *Task) Push(seeder Seeder) error {
	e := t.ct.CallTask(t, seeder)
	return e
}

// Call ...
func (t *Task) Call(seeder Seeder) error {
	return t.ct.CallTask(t, seeder)
}

// NewTask ...
func NewTask(task TaskAble) *Task {
	tsk := new(Task)
	tsk.ct = task
	return tsk
}
