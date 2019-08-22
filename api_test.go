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
	file, e := os.Open("D:\\videoall\\videos\\MIAA-086.wmv")
	if e != nil {
		t.Error(e)
		return
	}
	resolved, e := api.Unixfs().Add(context.Background(), files.NewReaderFile(file), func(settings *options.UnixfsAddSettings) error {
		settings.Pin = true
		return nil
	})
	if e != nil {
		t.Error(e)
	}
	t.Log(resolved)
}
