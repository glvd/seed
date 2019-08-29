package seed

// Task ...
type Task struct {
	Threader
}

// NewTask ...
func NewTask() Threader {
	tsk := new(Task)
	tsk.Threader = NewThread()
	return tsk
}

// TaskCall ...
func TaskCall(seeder Seeder, task *Task) error {
	return seeder.PushTo(StepperProcess, task)
}

// InformationTask ...
func InformationTask(task *Information) *Task {
	return &Task{
		Name:          "information",
		ProcessCaller: task,
	}
}
