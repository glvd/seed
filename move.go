package seed

import "context"

type move struct {
	to    string
	moves map[string]string
}

// Run ...
func (m *move) Run(context.Context) {

}

// BeforeRun ...
func (m *move) BeforeRun(seed *Seed) {

}

// AfterRun ...
func (m *move) AfterRun(seed *Seed) {
	panic("implement me")
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
