package seed

import (
	"context"

	httpapi "github.com/ipfs/go-ipfs-http-client"
)

type API struct {
	api *httpapi.HttpApi
	cb  chan APICallbackAble
}

func NewAPI() *API {
	return new(API)
}

type APICallbackAble interface {
	Callback(api *httpapi.HttpApi)
}

func (api *API) PushCallback(cb APICallbackAble) {
	api.cb <- cb
}

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
