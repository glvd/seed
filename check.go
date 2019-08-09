package seed

import (
	"context"

	"github.com/glvd/seed/model"
	httpapi "github.com/ipfs/go-ipfs-http-client"
	"github.com/ipfs/interface-go-ipfs-core/options"
)

// CheckType ...
type CheckType string

// CheckTypePin ...
const CheckTypePin CheckType = "pin"

// CheckTypeUnpin ...
const CheckTypeUnpin CheckType = "unpin"

// check ...
type check struct {
	api       *httpapi.HttpApi
	myID      *PeerID
	Type      string
	checkType CheckType
	from      []string
}

// Run ...
func (c *check) Run(context.Context) {
	log.Info("check running")
	switch c.checkType {
	case CheckTypePin:
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
	case CheckTypeUnpin:
		_, e := c.api.Pin().Ls(context.Background(), func(settings *options.PinLsSettings) error {
			settings.Type = c.Type
			return nil
		})
		if e != nil {
			log.Error(e)
			return
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

// CheckArgs ...
type CheckArgs func(c *check)

// CheckTypeArg ...
func CheckTypeArg(t string) CheckArgs {
	return func(c *check) {
		c.Type = t
	}
}

// CheckFromArg ...
func CheckFromArg(from ...string) CheckArgs {
	return func(c *check) {
		c.from = from
	}
}

// Check ...
func Check(args ...CheckArgs) Options {
	check := new(check)

	for _, argFn := range args {
		argFn(check)
	}
	return checkOption(check)
}

func checkOption(c *check) Options {
	return func(seed *Seed) {
		seed.thread[StepperCheck] = c
	}
}
