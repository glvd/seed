package task

import (
	"os"

	"github.com/glvd/seed"
)

// VideoProcess ...
type VideoProcess struct {
	Path  string
	Scale int64
	Skip  []interface{}
}

// CallTask ...
func (v *VideoProcess) CallTask(seeder seed.Seeder, task *seed.Task) error {
	select {
	case <-seeder.Context().Done():
		return nil
	default:
		files := seed.GetFiles(v.Path)
		for _, f := range files {
			if seed.IsVideo(f) {
				call := videoCall{
					path:  f,
					scale: v.Scale,
					skip:  v.Skip,
				}
				e := seeder.PushTo(seed.StepperProcess, call)
				if e != nil {
					log.Error(e)
					continue
				}
			}
		}
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

type videoCall struct {
	path  string
	scale int64
	skip  []interface{}
}

// Call ...
func (v *videoCall) Call(process *seed.Process) error {

}
