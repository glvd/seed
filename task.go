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
