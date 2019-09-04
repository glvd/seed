package task

import "os"

// VideoProcess ...
type VideoProcess struct {
	Path string
}

// NewVideoProcess ...
func NewVideoProcess() *VideoProcess {
	path := os.TempDir()
	return &VideoProcess{
		Path: path,
	}
}
