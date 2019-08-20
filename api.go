package seed

import (
	"context"

	httpapi "github.com/ipfs/go-ipfs-http-client"
	"github.com/ipfs/interface-go-ipfs-core/path"
	"github.com/multiformats/go-multiaddr"
)

// API ...
type API struct {
	api *httpapi.HttpApi
	cb  chan interface{}
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

// APICallbackAble ...
type APICallbackAble interface {
	Callback(*API, *httpapi.HttpApi) error
}

// CallbackFunc ...
type CallbackFunc func(*API, *httpapi.HttpApi) error

type cb struct {
	fn CallbackFunc
}

// PushCallback ...
func (api *API) PushCallback(cb APICallbackAble) {
	api.cb <- cb
}

// PushRun ...
func (api *API) PushRun(callbackFunc CallbackFunc) {
	api.cb <- &cb{
		fn: callbackFunc,
	}
}

// Run ...
func (api *API) Run(ctx context.Context) {
	log.Info("api running")
	var e error
	for {
		select {
		case <-ctx.Done():
			log.Info("api done")
		case c := <-api.cb:
			if v, b := c.(APICallbackAble); b {
				e = v.Callback(api, api.api)
				if e != nil {
					log.Error(e)
					continue
				}
			}
			if v, b := c.(APICallbackStatusAble); b {
				if e != nil {
					v.Failed()
					continue
				}
				v.Done()
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
func APIPeerID(api *API) *PeerID {
	pid := new(apiPeerID)
	pid.done = make(chan bool)
	go api.PushCallback(pid)
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
func (p *apiPeerID) Callback(api *httpapi.HttpApi) (e error) {
	p.id = new(PeerID)
	e = api.Request("id").Exec(context.Background(), p.id)
	if e != nil {
		return e
	}
	return nil
}

// APIPin ...
func APIPin(api *API, hash string) bool {
	p := new(apiPin)
	p.done = make(chan bool)
	go api.PushCallback(p)
	return <-p.done
}

type apiPin struct {
	hash string
	done chan bool
}

// Callback ...
func (a *apiPin) Callback(api *httpapi.HttpApi) error {
	return api.Pin().Add(context.Background(), path.New(a.hash))
}

// Done ...
func (a *apiPin) Done() {
	a.done <- true
}

// Failed ...
func (a *apiPin) Failed() {
	a.done <- false
}
