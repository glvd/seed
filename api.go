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
	*Thread
	failed *atomic.Bool
	api    *httpapi.HttpApi
	cb     chan APICaller
}

// Done ...
func (api *API) Done() <-chan bool {
	go func() {
		api.cb <- nil
	}()
	return api.Thread.Done()
}

// Option ...
func (api *API) Option(s Seeder) {
	apiOption(api)(s)
}

func apiOption(api *API) Options {
	return func(seeder Seeder) {
		seeder.SetBaseThread(StepperAPI, api)
	}
}

// IsFailed ...
func (api *API) IsFailed() bool {
	return api.failed.Load()
}

// Push ...
func (api *API) Push(v interface{}) error {
	return api.push(v)
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
	a.failed = atomic.NewBool(false)
	a.cb = make(chan APICaller, 10)
	a.Thread = NewThread()

	return a
}

// PushCallback ...
func (api *API) push(cb interface{}) (e error) {
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
			return
		case c := <-api.cb:
			if c == nil {
				break APIEnd
			}
			api.SetState(StateRunning)
			e = c.Call(api, api.api)
			if e != nil {
				log.Error(e)
			}
		case <-time.After(30 * time.Second):
			log.Info("api time out")
			api.SetState(StateWaiting)
		}
	}
	close(api.cb)
	api.Finished()
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

var _ APICaller = &apiCall{}
