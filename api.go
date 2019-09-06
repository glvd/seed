package seed

import (
	"context"
	"errors"
	"os"
	"time"

	files "github.com/ipfs/go-ipfs-files"
	httpapi "github.com/ipfs/go-ipfs-http-client"
	"github.com/ipfs/interface-go-ipfs-core/options"
	"github.com/ipfs/interface-go-ipfs-core/path"
	"github.com/multiformats/go-multiaddr"
	"go.uber.org/atomic"
)

// API ...
type API struct {
	*Thread
	failed *atomic.Bool
	api    *httpapi.HttpApi
	cb     chan APICaller
}

// Failed ...
func (api *API) Failed() bool {
	return api.failed.Load()
}

// SetFailed ...
func (api *API) SetFailed(failed bool) {
	api.failed.Store(failed)
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
	return api.push(v)
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
	a.failed = atomic.NewBool(false)
	a.cb = make(chan APICaller, 10)
	a.Thread = NewThread()

	return a
}

// PushCallback ...
func (api *API) push(cb interface{}) (e error) {
	if v, b := cb.(APICaller); b {
		api.cb <- v
		return nil
	}
	return errors.New("not api callback")
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
				break APIEnd
			}
			api.SetState(StateRunning)
			e = c.Call(api, api.api)
			if e != nil {
				log.Error(e)
			}
		case <-time.After(30 * time.Second):
			log.Info("api time out")
			api.SetState(StateWaiting)
		}
	}
	close(api.cb)
	api.Finished()
}

// PeerID ...
type PeerID struct {
	Addresses       []string `json:"Addresses"`
	AgentVersion    string   `json:"AgentVersion"`
	ID              string   `json:"ID"`
	ProtocolVersion string   `json:"ProtocolVersion"`
	PublicKey       string   `json:"PublicKey"`
}

// APICallback ...
func APICallback(v interface{}, cb APICallbackFunc) (Stepper, APICaller) {
	return StepperAPI, &apiCall{
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

var _ APICaller = &apiCall{}

// AddFile ...
func AddFile(api *API, filename string) (path.Resolved, error) {
	file, e := os.Open(filename)
	if e != nil {
		return nil, e
	}
	resolved, e := api.api.Unixfs().Add(api.Context(), files.NewReaderFile(file),
		func(settings *options.UnixfsAddSettings) error {
			settings.Pin = true
			return nil
		})
	return resolved, e
}

// AddDir ...
func AddDir(api *API, dir string) (path.Resolved, error) {
	stat, err := os.Lstat(dir)
	if err != nil {
		return nil, err
	}

	sf, err := files.NewSerialFile(dir, false, stat)
	if err != nil {
		return nil, err
	}
	//不加目录
	//slf := files.NewSliceDirectory([]files.DirEntry{files.FileEntry(filepath.Base(dir), sf)})
	//reader := files.NewMultiFileReader(slf, true)
	resolved, e := api.api.Unixfs().Add(api.Context(), sf,
		func(settings *options.UnixfsAddSettings) error {
			settings.Pin = true
			return nil
		})
	return resolved, e
}
