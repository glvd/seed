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
	for _, v := range m.moves {
		_, name := filepath.Split(v)
		e = os.Rename(v, filepath.Join(m.to, name))
		if e != nil {
			log.Error(e)
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
	//output nothing
}

// Move ...
func Move(path string) Options {
	info := &move{
		to: path,
	}
	return moveOption(info)
}

func moveOption(move *move) Options {
	return func(seed *Seed) {
		seed.thread[StepperMove] = move
	}
}
