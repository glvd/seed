package seed

import (
	"context"
)

// Slice ...
type Slice struct {
	Seed  *seed
	Scale int64
	File  chan string
}

// Run ...
func (s *Slice) Run(context.Context) {

}

// BeforeRun ...
func (s *Slice) BeforeRun(seed *seed) {
	panic("implement me")
}

// AfterRun ...
func (s *Slice) AfterRun(seed *seed) {
	panic("implement me")
}
