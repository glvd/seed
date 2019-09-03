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
const CheckTypePin CheckType = "pin"

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

// CallTask ...
func (c *Check) CallTask(seeder Seeder, task *Task) error {
	pins := make(chan []iface.Pin)
	select {
	case <-seeder.Context().Done():
		return nil
	default:
		e := seeder.PushTo(APICallback(c.MyID, func(api *API, api2 *httpapi.HttpApi, v interface{}) (e error) {
			pid := v.(*PeerID)
			e = api2.Request("id").Exec(seeder.Context(), pid)
			if e != nil {
				return e
			}
			plist, err := api2.Pin().Ls(seeder.Context(), func(settings *options.PinLsSettings) error {
				settings.Type = c.Type
				return nil
			})
			if err != nil {
				return err
			}
			pins <- plist
			return nil
		}))
		if e != nil {
			return e
		}
		//waiting for result
		for _, value := range <-pins {
			p := &model.Pin{
				PinHash: model.PinHash(value.Path()),
				PeerID:  []string{c.MyID.ID},
				VideoID: "",
			}
			e = seeder.PushTo(DatabaseCallback(p, func(database *Database, eng *xorm.Engine, v interface{}) (e error) {
				return nil
			}))
			if e != nil {
				log.Error(e)
			}
		}
	}

	return nil
}

// Run ...
func (c *Check) Run(ctx context.Context) {

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
