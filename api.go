package seed

import httpapi "github.com/ipfs/go-ipfs-http-client"

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
