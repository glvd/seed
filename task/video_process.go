package task

import (
	"os"

	"github.com/glvd/seed"
)

// VideoProcess ...
type VideoProcess struct {
	Path string
}

// CallTask ...
func (v *VideoProcess) CallTask(seeder seed.Seeder, task *seed.Task) error {
	select {
	case <-seeder.Context().Done():
		return nil
	default:

	}

	return nil
}

// NewVideoProcess ...
func NewVideoProcess() *VideoProcess {
	path := os.TempDir()
	return &VideoProcess{
		Path: path,
	}
}

// Task ...
func (v *VideoProcess) Task() *seed.Task {
	return seed.NewTask(v)
}
