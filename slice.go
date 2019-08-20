package seed

import (
	"context"

	cmd "github.com/godcong/go-ffmpeg-cmd"
)

// Slice ...
type Slice struct {
	Seed  *Seed
	Scale int64
	File  chan string
}

// Run ...
func (s *Slice) Run(context.Context) {

}

// BeforeRun ...
func (s *Slice) BeforeRun(seed *Seed) {
	panic("implement me")
}

// AfterRun ...
func (s *Slice) AfterRun(seed *Seed) {
	panic("implement me")
}
