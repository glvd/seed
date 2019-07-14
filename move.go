package seed

import (
	"context"
	"os"
	"path/filepath"
)

type move struct {
	to    string
	moves map[string]string
}

// Run ...
func (m *move) Run(context.Context) {
	var e error
	s, e := filepath.Abs(m.to)
	if e != nil {
		return
	}
	for hash, v := range m.moves {
		//_, name := filepath.Split(v)
		to := hash + filepath.Ext(v)
		path := filepath.Join(s, to)
		log.With("from", path, "to", to).Info("move")
		e = os.Rename(v, path)
		if e != nil {
			log.Error(e, path)
			continue
		}
	}
}

// BeforeRun ...
func (m *move) BeforeRun(seed *Seed) {
	m.moves = seed.Moves
}

// AfterRun ...
func (m *move) AfterRun(seed *Seed) {
	seed.Moves = make(map[string]string)
}

// MoveInfo ...
func MoveInfo(path string) Options {
	info := &move{
		to: path,
	}
	return moveOption(StepperMoveInfo, info)
}

// MoveProc ...
func MoveProc(path string) Options {
	proc := &move{
		to: path,
	}
	return moveOption(StepperMoveproc, proc)
}

func moveOption(m Stepper, move *move) Options {
	return func(seed *Seed) {
		seed.thread[m] = move
	}
}
