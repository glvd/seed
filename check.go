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

// Check ...
type Check struct {
	*Seed
	api       *httpapi.HttpApi
	myID      *PeerID
	Type      string
	checkType CheckType
	from      []string
	skipType  []interface{}
}

// Option ...
func (c *Check) Option(seed *Seed) {
	checkOption(c)(seed)
}

// Run ...
func (c *Check) Run(context.Context) {
	log.Info("Check running")
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
			e := model.UpdateVideo()
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
				s.NotIn("type", c.skipType...)
			}
			i, e := s.Clone().Count(model.Unfinished{})
			if e != nil {
				log.Error(e)
				return
			}
			for start := 0; start < int(i); start += 50 {
				unfins, e := model.AllUnfinished(s.Clone(), 50, start)
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
func (c *Check) BeforeRun(seed *Seed) {
	c.myID = APIPeerID(seed.API)
	if c.Type == "" {
		c.Type = "recursive"
	}
}

// AfterRun ...
func (c *Check) AfterRun(seed *Seed) {

}

// CheckArgs ...
type CheckArgs func(c *Check)

// CheckSkipArg ...
func CheckSkipArg(s ...string) CheckArgs {
	return func(c *Check) {
		for i := range s {
			c.skipType = append(c.skipType, s[i])
		}
	}
}

// CheckTypeArg ...
func CheckTypeArg(t CheckType) CheckArgs {
	return func(c *Check) {
		c.checkType = t
	}
}

// CheckPinTypeArg ...
func CheckPinTypeArg(t string) CheckArgs {
	return func(c *Check) {
		c.Type = t
	}
}

// CheckFromArg ...
func CheckFromArg(from ...string) CheckArgs {
	return func(c *Check) {
		c.from = from
	}
}

// NewCheck ...
func NewCheck(args ...CheckArgs) *Check {
	check := new(Check)

	for _, argFn := range args {
		argFn(check)
	}
	return check
}

func checkOption(c *Check) Options {
	return func(seed *Seed) {
		seed.thread[StepperCheck] = c
	}
}
