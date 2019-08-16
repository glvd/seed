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
	Callback(api *httpapi.HttpApi)
}

// PushCallback ...
func (api *API) PushCallback(cb APICallbackAble) {
	api.cb <- cb
}

// Run ...
func (api *API) Run(ctx context.Context) {
	log.Info("api running")
	for {
		select {
		case <-ctx.Done():
			log.Info("api done")
		case c := <-api.cb:
			c.Callback(api.api)
		}
	}
}
