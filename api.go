package seed

import (
	"context"
	"sync"

	httpapi "github.com/ipfs/go-ipfs-http-client"
	"github.com/ipfs/interface-go-ipfs-core/path"
	"github.com/multiformats/go-multiaddr"
	"golang.org/x/xerrors"
)

// API ...
type API struct {
	Seeder
	api  *httpapi.HttpApi
	cb   chan APICaller
	done chan bool
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
		seeder.SetBaseThread(StepperAPI, api)
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
	return a
}

// CallbackFunc ...
type CallbackFunc func(*API, *httpapi.HttpApi) error

// PushCallback ...
func (api *API) pushAPICallback(cb interface{}) (e error) {
	if v, b := cb.(APICaller); b {
		//go func(c APICaller) {
		api.cb <- v
		//}(v)
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
				log.Info("api end")
				break APIEnd
			}
			e = c.Call(api, api.api)
			if e != nil {
				log.Error(e)
			}
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

// APIPeerID ...
func APIPeerID(seed Seeder) *PeerID {
	pid := new(apiPeerID)
	pid.done = make(chan bool)
	e := seed.PushTo(StepperAPI, pid)
	if e != nil {
		return nil
	}
	d := <-pid.done
	if d {
		return pid.id
	}
	return nil
}

type apiPeerID struct {
	id   *PeerID
	done chan bool
}

// Done ...
func (p *apiPeerID) Done() {
	p.done <- true
}

// Failed ...
func (p *apiPeerID) Failed() {
	p.done <- false
}

// OnDone ...
func (p *apiPeerID) OnDone() *PeerID {
	d := <-p.done
	if d {
		return p.id
	}

	return nil
}

// Callback ...
func (p *apiPeerID) Callback(api *API, api2 *httpapi.HttpApi) (e error) {
	p.id = new(PeerID)
	e = api2.Request("id").Exec(context.Background(), p.id)
	if e != nil {
		return e
	}
	return nil
}

// APIPin ...
func APIPin(seed Seeder, hash string) (e error) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	e = seed.PushTo(StepperAPI, func(api *API, api2 *httpapi.HttpApi) error {
		defer wg.Done()
		e = api2.Pin().Add(context.Background(), path.New(hash))
		return e
	})
	wg.Wait()
	return e
}

// APICallback ...
func APICallback(v interface{}, cb APICallbackFunc) APICaller {
	return &apiCall{
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

// APIOption ...
func APIOption(s string) Options {
	return func(seed Seeder) {
		seed.SetThread(StepperAPI, NewAPI(s))
	}
}
