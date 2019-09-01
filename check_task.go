package seed

import (
	"context"

	"github.com/glvd/seed/model"
	"github.com/go-xorm/xorm"
	httpapi "github.com/ipfs/go-ipfs-http-client"
	iface "github.com/ipfs/interface-go-ipfs-core"
	"github.com/ipfs/interface-go-ipfs-core/options"
)

// CheckType ...
type CheckType string

// CheckTypePin ...
const CheckTypePin CheckType = "Pin"

// CheckTypeUnpin ...
const CheckTypeUnpin CheckType = "unpin"

// Check ...
type Check struct {
	MyID      *PeerID
	Type      string
	CheckType CheckType
	from      []string
	skipType  []interface{}
}

func (c *Check) CallTask(seeder Seeder, taks *Task) error {
	select {
	case <-seeder.Context().Done():
		return nil
	default:
		e := seeder.PushTo(APICallback(c.MyID, func(api *API, api2 *httpapi.HttpApi, v interface{}) (e error) {
			return nil
		}))
		if e != nil {
			return e
		}

	}

	return nil
}

// Run ...
func (c *Check) Run(ctx context.Context) {
	log.Info("Check running")
	cPin := make(chan iface.Pin)
	e := c.PushTo(StepperAPI, func(api *API, api2 *httpapi.HttpApi) (e error) {
		defer func() {
			cPin <- nil
		}()
		pins, e := api2.Pin().Ls(ctx, func(settings *options.PinLsSettings) error {
			settings.Type = c.Type
			return nil
		})
		if e != nil {
			return e
		}
		for _, p := range pins {
			cPin <- p
		}

		return nil
	})
	if e != nil {
		return
	}
	switch c.checkType {
	case CheckTypePin:
	PinList:
		for {
			select {
			case <-ctx.Done():
				return
			case pin := <-cPin:
				if pin == nil {
					break PinList
				}
				log.With("path", pin.Path()).Info("pinned")
				p := &model.Pin{
					PinHash: model.PinHash(pin.Path()),
					PeerID:  []string{c.myID.ID},
					VideoID: "",
				}
				e := c.PushTo(StepperDatabase, func(database *Database, eng *xorm.Engine) (e error) {
					return model.UpdatePinVideoID(eng.NewSession(), p)
				})
				if e != nil {
					return
				}

			}
		}

	case CheckTypeUnpin:
		unf := make(chan *model.Unfinished)
		c.PushTo(StepperDatabase, func(database *Database, eng *xorm.Engine) (e error) {
			defer func() {
				unf <- nil
			}()
			session := eng.NewSession()
			if len(c.skipType) > 0 {
				session.NotIn("type", c.skipType...)
			}
			i, e := session.Clone().Count(model.Unfinished{})
			if e != nil {
				return e
			}
			for start := 0; start < int(i); start += 50 {
				unfins, e := model.AllUnfinished(session.Clone(), 50, start)
				if e != nil {
					log.Error(e)
					return e
				}

				log.Infof("Pin(%d)", len(*unfins))
				for i := range *unfins {
					unf <- (*unfins)[i]
				}
			}
			return nil
		})

		//var retUnf []*model.Unfinished
		pins := make(map[string][]byte)
	PinList2:
		for {
			select {
			case <-ctx.Done():
				return
			case pin := <-cPin:
				if pin == nil {
					break PinList2
				}
				pins[model.PinHash(pin.Path())] = nil
			}
		}
	CheckList2:
		for {
			select {
			case u := <-unf:
				if u == nil {
					break CheckList2
				}
				if _, b := pins[u.Hash]; b {
					log.With("hash", u.Hash, "relate", u.Relate, "type", u.Type).Info("unpin")
				}
			}
		}

	}

}

// BeforeRun ...
func (c *Check) BeforeRun(seed Seeder) {
	//c.myID = APIPeerID(seed)
	if c.Type == "" {
		c.Type = "recursive"
	}
}

// AfterRun ...
func (c *Check) AfterRun(seed Seeder) {

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
	return func(seed Seeder) {
		seed.SetThread(StepperCheck, c)
	}
}
