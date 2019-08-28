package seed

import (
	"context"
	"time"
)

// MoveCaller ...
type MoveCaller interface {
	Call(*Move) error
}

// Move ...
type Move struct {
	*Thread
	cb    chan MoveCaller
	to    string
	Moves map[string]string
}

// Register ...
func (m *Move) Register(ops ...Optioner) {
	panic("implement me")
}

// Push ...
func (m *Move) Push(interface{}) error {
	return nil
}

// NewMove ...
func NewMove() *Move {
	return &Move{
		Thread: NewThread(),
	}
}

// Run ...
func (m *Move) Run(ctx context.Context) {
	log.Info("move running")

InfoEnd:
	for {
		select {
		case <-ctx.Done():
			break InfoEnd
		case cb := <-m.cb:
			if cb == nil {
				break InfoEnd
			}
			m.SetState(StateRunning)
			e := cb.Call(m)
			if e != nil {
				log.Error(e)
			}
		case <-time.After(30 * time.Second):
			log.Info("info time out")
			m.SetState(StateWaiting)
		}
	}
	close(m.cb)
	m.Finished()
}

// MoveInfo ...
func MoveInfo(path string) Options {
	info := &Move{
		to: path,
	}
	return MoveOption(StepperMoveInfo, info)
}

// MoveProc ...
func MoveProc(path string) Options {
	proc := &Move{
		to: path,
	}
	return MoveOption(StepperMoveproc, proc)
}

// MoveOption ...
func MoveOption(stepper Stepper, Move *Move) Options {
	return func(seed Seeder) {
		seed.SetThread(stepper, Move)
	}
}
