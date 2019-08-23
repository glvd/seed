package seed

import (
	"context"

	"go.uber.org/atomic"
)

// BaseThread ...
type BaseThread struct {
	Seeder
	state *atomic.Int32
	done  chan bool
}

// Finished ...
func (t *BaseThread) Finished() {
	t.done <- true
}

// Run ...
func (t *BaseThread) Run(context.Context) {
	panic("implement me")
}

// SetState ...
func (t *BaseThread) SetState(state State) {
	t.state.Store(int32(state))
}

// Push ...
func (t *BaseThread) Push(interface{}) error {
	panic("implement me")
}

// BeforeRun ...
func (t *BaseThread) BeforeRun(seed Seeder) {
	t.Seeder = seed
}

// AfterRun ...
func (t *BaseThread) AfterRun(seed Seeder) {
}

// State ...
func (t *BaseThread) State() State {
	return State(t.state.Load())
}

// Done ...
func (t *BaseThread) Done() <-chan bool {
	return t.done
}

// NewThread ...
func NewThread() Threader {
	return &BaseThread{
		state: atomic.NewInt32(int32(StateWaiting)),
		done:  make(chan bool),
	}
}
