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
	//Step() Stepper
	CallTask(Seeder, *Task) error
}

// Task ...
type Task struct {
	ct       TaskAble
	taskStep TaskStep
	//*Thread
}

// Push ...
func (t *Task) Push(seeder Seeder) error {
	e := t.ct.CallTask(seeder, t)
	return e
}

// NewTask ...
func NewTask(task TaskAble) *Task {
	tsk := new(Task)
	tsk.ct = task
	return tsk
}
