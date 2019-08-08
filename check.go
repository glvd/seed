package seed

import (
	"context"

	shell "github.com/godcong/go-ipfs-restapi"
)

// check ...
type check struct {
	shell *shell.Shell
	Type  string
}

func (c *check) Run(context.Context) {
	switch c.Type {
	case "pin":

	}
}

// BeforeRun ...
func (c *check) BeforeRun(seed *Seed) {
	c.shell = seed.Shell
}

// AfterRun ...
func (c *check) AfterRun(seed *Seed) {

}

// Process ...
func Check(tp string) Options {
	check := &check{
		Type: tp,
	}
	return checkOption(check)
}

func checkOption(c *check) Options {
	return func(seed *Seed) {
		seed.thread[StepperCheck] = c
	}
}
