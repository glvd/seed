package seed

import (
	"context"

	httpapi "github.com/ipfs/go-ipfs-http-client"
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

// MyPeerID ...
func (api *API) MyPeerID() (pid *PeerID, e error) {
	pid = new(PeerID)
	e = api.api.Request("id").Exec(context.Background(), pid)
	return
}

type peerID struct {
	id   *PeerID
	done chan bool
}

// Callback ...
func (p *peerID) Callback(api *httpapi.HttpApi) (e error) {
	p.id = new(PeerID)
	e = api.Request("id").Exec(context.Background(), p.id)
	if e != nil {
		return e
	}
	return nil
}
