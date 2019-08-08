package seed

import (
	"context"

	"github.com/glvd/seed/model"
	httpapi "github.com/ipfs/go-ipfs-http-client"
	"github.com/ipfs/interface-go-ipfs-core/options"
)

// check ...
type check struct {
	api  *httpapi.HttpApi
	myID *PeerID
	Type string
}

func (c *check) Run(context.Context) {
	log.Info("check running")
	pins, e := c.api.Pin().Ls(context.Background(), func(settings *options.PinLsSettings) error {
		settings.Type = c.Type
		return nil
	})
	if e != nil {
		log.Error(e)
		return
	}
	for _, path := range pins {
		log.With("path", path.Path()).Info("pinned")
		p := &model.Pin{
			PinHash: model.PinHash(path.Path()),
			PeerID:  []string{c.myID.ID},
			VideoID: "",
		}
		e := p.UpdateVideo()
		if e != nil {
			log.Error(e)
		}
	}

}

// BeforeRun ...
func (c *check) BeforeRun(seed *Seed) {
	var e error
	c.api = seed.API
	c.myID, e = seed.MyPeerID()
	if e != nil {
		log.Error(e)
	}
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
