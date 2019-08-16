package seed

import (
	"context"

	httpapi "github.com/ipfs/go-ipfs-http-client"
	"github.com/ipfs/interface-go-ipfs-core/path"
)

// API ...
type API struct {
	api *httpapi.HttpApi
	cb  chan APICallbackAble
}

// NewAPI ...
func NewAPI() *API {
	return new(API)
}

// APICallbackAble ...
type APICallbackAble interface {
	Callback(api *httpapi.HttpApi) error
	Done()
	Failed()
}

// PushCallback ...
func (api *API) PushCallback(cb APICallbackAble) {
	api.cb <- cb
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
			e = c.Callback(api.api)
			if e != nil {
				log.Error(e)
				c.Failed()
				continue
			}
			c.Done()
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

// MyPeerID ...
func MyPeerID(api *API) *PeerID {
	pid := new(apiPeerID)
	go api.PushCallback(pid)
	return pid.OnDone()
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
	select {
	case d := <-p.done:
		if d {
			return p.id
		}
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

}

// Failed ...
func (a *apiPin) Failed() {

}

// OnDone ...
func (a *apiPin) OnDone() {

}
