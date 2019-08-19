package seed

import (
	"context"
	"os"
	"path/filepath"
)

// Move ...
type Move struct {
	to    string
	Moves map[string]string
}

// NewMove ...
func NewMove() *Move {
	return &Move{}
}

// Run ...
func (m *Move) Run(context.Context) {
	var e error
	s, e := filepath.Abs(m.to)
	if e != nil {
		return
	}
	for v, hash := range m.Moves {
		//_, name := filepath.Split(v)
		to := hash + filepath.Ext(v)
		path := filepath.Join(s, to)
		log.With("from", v, "to", to).Info("Move")
		e = os.Rename(v, path)
		if e != nil {
			log.Error(e, path)
			continue
		}
	}
}

// BeforeRun ...
func (m *Move) BeforeRun(seed *Seed) {
}

// AfterRun ...
func (m *Move) AfterRun(seed *Seed) {
	seed.Moves = make(map[string]string)
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
func MoveOption(m Stepper, Move *Move) Options {
	return func(seed *Seed) {
		seed.thread[m] = Move
	}
}
