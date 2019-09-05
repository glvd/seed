package seed_test

import (
	"context"
	"os"
	"testing"

	files "github.com/ipfs/go-ipfs-files"
	httpapi "github.com/ipfs/go-ipfs-http-client"
	"github.com/ipfs/interface-go-ipfs-core/options"
	"github.com/multiformats/go-multiaddr"
)

// TestApiCall_Call ...
func TestApiAdd(t *testing.T) {
	var api *httpapi.HttpApi
	addr, e := multiaddr.NewMultiaddr("/ip4/127.0.0.1/tcp/5001")
	if e != nil {
		t.Error(e)
		return
	}
	api, e = httpapi.NewApi(addr)
	if e != nil {
		t.Error(e)
		return
	}
	file, e := os.Open("d:\\video\\temp\\9df762a1-b1ae-478f-82c6-462d7b3a2286\\")
	if e != nil {
		t.Error(e)
		return
	}
	np, e := files.NewReaderPathFile("d:\\video\\temp\\9df762a1-b1ae-478f-82c6-462d7b3a2286\\", file, nil)
	if e != nil {
		t.Error(e)
		return
	}
	resolved, e := api.Unixfs().Add(context.Background(), np, func(settings *options.UnixfsAddSettings) error {
		settings.Pin = true

		return nil
	})
	if e != nil {
		t.Error(e)
	}
	t.Log(resolved)
}
