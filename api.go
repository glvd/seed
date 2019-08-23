package seed

import (
	"context"
	"time"

	httpapi "github.com/ipfs/go-ipfs-http-client"
	"github.com/multiformats/go-multiaddr"
	"go.uber.org/atomic"
	"golang.org/x/xerrors"
)

// API ...
type API struct {
	Seeder
	api   *httpapi.HttpApi
	cb    chan APICaller
	done  chan bool
	state *atomic.Int32
}

// State ...
func (api *API) State() State {
	return State(api.state.Load())
}

// Done ...
func (api *API) Done() <-chan bool {
	go func() {
		api.cb <- nil
	}()
	return api.done
}

// Option ...
func (api *API) Option(s Seeder) {
	apiOption(api)(s)
}

func apiOption(api *API) Options {
	return func(seeder Seeder) {
		seeder.SetThread(StepperAPI, api)
	}
}

// Push ...
func (api *API) Push(v interface{}) error {
	return api.pushAPICallback(v)
}

// BeforeRun ...
func (api *API) BeforeRun(seed Seeder) {
	api.Seeder = seed
}

// AfterRun ...
func (api *API) AfterRun(seed Seeder) {
}

// NewAPI ...
func NewAPI(path string) *API {
	a := new(API)
	var e error
	addr, e := multiaddr.NewMultiaddr(path)
	if e != nil {
		panic(e)
	}
	a.api, e = httpapi.NewApi(addr)
	if e != nil {
		panic(e)
	}
	a.cb = make(chan APICaller, 10)
	a.done = make(chan bool)
	a.state = atomic.NewInt32(int32(StateWaiting))
	return a
}

// PushCallback ...
func (api *API) pushAPICallback(cb interface{}) (e error) {
	if v, b := cb.(APICaller); b {
		api.cb <- v
		return
	}
	return xerrors.New("not api callback")
}

// Run ...
func (api *API) Run(ctx context.Context) {
	log.Info("api running")
	var e error
APIEnd:
	for {
		select {
		case <-ctx.Done():
			api.state.Store(int32(StateStop))
			return
		case c := <-api.cb:
			if c == nil {
				api.state.Store(int32(StateStop))
				break APIEnd
			}
			api.state.Store(int32(StateRunning))
			e = c.Call(api, api.api)
			if e != nil {
				log.Error(e)
			}
		case <-time.After(30 * time.Second):
			api.state.Store(int32(StateWaiting))
		}
	}
	close(api.cb)
	api.done <- true
}

// PeerID ...
type PeerID struct {
	Addresses       []string `json:"Addresses"`
	AgentVersion    string   `json:"AgentVersion"`
	ID              string   `json:"ID"`
	ProtocolVersion string   `json:"ProtocolVersion"`
	PublicKey       string   `json:"PublicKey"`
}

// APICallback ...
func APICallback(v interface{}, cb APICallbackFunc) (Stepper, APICaller) {
	return StepperAPI, &apiCall{
		v:  v,
		cb: cb,
	}
}

type apiCall struct {
	v  interface{}
	cb APICallbackFunc
}

// Callback ...
func (a *apiCall) Call(api *API, api2 *httpapi.HttpApi) error {
	return a.cb(api, api2, a.v)
}
