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
	Seed *seed
	api  *httpapi.HttpApi
	cb   chan APICallback
}

// Push ...
func (api *API) Push(v interface{}) error {
	return api.pushAPICallback(v)
}

// BeforeRun ...
func (api *API) BeforeRun(seed *seed) {
	api.Seed = seed
}

// AfterRun ...
func (api *API) AfterRun(seed *seed) {
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
	return a
}

// APICallbackStatusAble ...
type APICallbackStatusAble interface {
	Done()
	Failed()
}

// APICallback ...
type APICallback func(*API, *httpapi.HttpApi) error

// APICallbackAble ...
type APICallbackAble interface {
	Callback(*API, *httpapi.HttpApi) error
}

// CallbackFunc ...
type CallbackFunc func(*API, *httpapi.HttpApi) error

type apiCallback struct {
	fn APICallback
}

// PushCallback ...
func (api *API) pushAPICallback(cb interface{}) (e error) {
	if v, b := cb.(APICallback); b {
		go func(a APICallback) {
			api.cb <- a
		}(v)
	}
	return xerrors.New("not api callback")
}

// Run ...
func (api *API) Run(ctx context.Context) {
	log.Info("api running")
	var e error
	for {
		select {
		case <-ctx.Done():
			return
		case c := <-api.cb:
			e = c(api, api.api)
			if e != nil {
				log.Error(e)
			}
		}
	}
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
func APIPeerID(seed *seed) *PeerID {
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
func APIPin(seed *seed, hash string) (e error) {
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

type apiPin struct {
	hash string
	done chan bool
}

// Callback ...
func (a *apiPin) Callback(api *API, api2 *httpapi.HttpApi) error {
	return api2.Pin().Add(context.Background(), path.New(a.hash))
}

// Done ...
func (a *apiPin) Done() {
	a.done <- true
}

// Failed ...
func (a *apiPin) Failed() {
	a.done <- false
}
