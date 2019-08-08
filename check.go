package seed

import (
	"context"

	httpapi "github.com/ipfs/go-ipfs-http-client"
	"github.com/ipfs/interface-go-ipfs-core/options"
)

// check ...
type check struct {
	api  *httpapi.HttpApi
	Type string
}

func (c *check) Run(context.Context) {
	pins, e := c.api.Pin().Ls(context.Background(), func(settings *options.PinLsSettings) error {
		settings.Type = c.Type
		return nil
	})
	if e != nil {
		log.Error(e)
		return
	}
	for _, p := range pins {
		log.With("path", p.Path()).Info("pinned")
	}

}

// BeforeRun ...
func (c *check) BeforeRun(seed *Seed) {
	c.api = seed.API
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
