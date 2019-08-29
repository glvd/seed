package seed

// Task ...
type Task struct {
	Name string
	ProcessCaller
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
