package seed

// TaskProcessor ...
type TaskProcessor interface {
	Call(process *Process) error
}

// Task ...
type Task struct {
	Name string
	TaskProcessor
}

// TaskCall ...
func TaskCall(seeder Seeder, task *Task) error {
	return seeder.PushTo(StepperProcess, task)
}

// InformationTask ...
func InformationTask(task *Information) *Task {
	return &Task{
		Name:          "information",
		TaskProcessor: task,
	}
}
