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
	skipType  []string
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
		pins, e := c.api.Pin().Ls(context.Background(), func(settings *options.PinLsSettings) error {
			settings.Type = c.Type
			return nil
		})
		if e != nil {
			log.Error(e)
			return
		}
		unf := make(chan *model.Unfinished)

		go func(u chan<- *model.Unfinished) {
			defer func() {
				u <- nil
			}()
			s := model.DB().NewSession()
			if len(c.skipType) > 0 {
				for idx := range c.skipType {
					s = s.Or("type = ?", c.skipType[idx])
				}
			}
			i, e := s.Clone().Count(model.Unfinished{})
			if e != nil {
				log.Error(e)
				return
			}
			log.Info(s.LastSQL())
			for start := 0; start < int(i); start += 50 {
				unfins, e := model.AllUnfinished(s, 50, start)
				if e != nil {
					log.Error(e)
					return
				}

				log.Infof("pin(%d)", len(*unfins))
				for i := range *unfins {
					u <- (*unfins)[i]
				}
			}

		}(unf)
		//var retUnf []*model.Unfinished
	CheckEnd:
		for {
			select {
			case u := <-unf:
				if u == nil {
					break CheckEnd
				}
				pinned := false
				for _, path := range pins {
					if u.Hash == model.PinHash(path.Path()) {
						pinned = true
					}
					//retUnf = append(retUnf, u)
				}
				if !pinned {
					log.With("hash", u.Hash, "relate", u.Relate, "type", u.Type).Info("unpin")
				}

			}
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
	if c.Type == "" {
		c.Type = "recursive"
	}
}

// AfterRun ...
func (c *check) AfterRun(seed *Seed) {

}

// CheckArgs ...
type CheckArgs func(c *check)

// CheckSkipArg ...
func CheckSkipArg(s []string) CheckArgs {
	return func(c *check) {
		c.skipType = s
	}
}

// CheckTypeArg ...
func CheckTypeArg(t CheckType) CheckArgs {
	return func(c *check) {
		c.checkType = t
	}
}

// CheckPinTypeArg ...
func CheckPinTypeArg(t string) CheckArgs {
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
