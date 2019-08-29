package seed

// TaskProcessor ...
type TaskProcessor interface {
	Call(process *Process)
}

// Task ...
type Task struct {
	Name string
	TaskProcessor
}
