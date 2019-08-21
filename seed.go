package seed

import (
	"context"
	"crypto/sha1"
	"fmt"
	"math"
	"sync"

	"github.com/glvd/seed/model"
	shell "github.com/godcong/go-ipfs-restapi"
	jsoniter "github.com/json-iterator/go"
	"golang.org/x/xerrors"
)

// Options ...
type Options func(seeder Seeder)

// AfterInitOptions ...
type AfterInitOptions func(Seeder)

// Thread ...
type Thread struct {
	wg sync.WaitGroup
}

// seed ...
type seed struct {
	Shell       *shell.Shell
	API         *API
	Move        *Move
	Database    *Database
	Workspace   string
	Scale       int64
	NoCheck     bool
	Unfinished  map[string]*model.Unfinished
	Videos      map[string]*model.Video
	Moves       map[string]string
	MaxLimit    int
	From        string
	args        map[string]interface{}
	wg          *sync.WaitGroup
	ctx         context.Context
	cancel      context.CancelFunc
	skipConvert bool
	preAdd      bool
	noSlice     bool
	upScale     bool
	threads     int
	thread      map[Stepper]Threader
	ignores     map[string][]byte
	err         error
	skipExist   bool
}

// GetThread ...
func (seed *seed) GetThread(stepper Stepper) Threader {
	return seed.thread[stepper]
}

// SetThread ...
func (seed *seed) SetThread(stepper Stepper, threader Threader) {
	seed.thread[stepper] = threader
}

// HasThread ...
func (seed *seed) HasThread(stepper Stepper) bool {
	_, b := seed.thread[stepper]
	return b
}

// PushTo ...
func (seed *seed) PushTo(stepper Stepper, v interface{}) (e error) {
	if val, b := seed.thread[stepper]; b {
		return val.Push(v)
	}
	return xerrors.Errorf("thread(%d) is not exist")
}

// Args ...
func (seed *seed) Args() map[string]interface{} {
	return seed.args
}

// SetArgs ...
func (seed *seed) SetArgs(args map[string]interface{}) {
	seed.args = args
}

// AddArg ...
func (seed *seed) AddArg(key string, value interface{}) {
	seed.args[key] = value
}

// GetArg ...
func (seed *seed) GetArg(key string) (v interface{}, b bool) {
	v, b = seed.args[key]
	return
}

// GetStringArg ...
func (seed *seed) GetStringArg(key string) (v string) {
	if arg, b := seed.GetArg(key); b {
		v, _ = arg.(string)
	}
	return
}

// GetBoolArg ...
func (seed *seed) GetBoolArg(key string) (v bool) {
	if arg, b := seed.GetArg(key); b {
		v, _ = arg.(bool)
	}
	return
}

// GetNumberArg ...
func (seed *seed) GetNumberArg(key string) (v int64) {
	if arg, b := seed.GetArg(key); b {
		v, _ = arg.(int64)
	}
	return
}

// Stop ...
func (seed *seed) Stop() {
	if seed.cancel != nil {
		seed.cancel()
	}
}

// Err ...
func (seed *seed) Err() error {
	return seed.err
}

// Start ...
func (seed *seed) Start() {
	log.Info("Seed starting")
	for i := range seed.thread {
		if seed.thread[i] == nil {
			continue
		}
		seed.wg.Add(1)
		go func(t Threader, group *sync.WaitGroup) {
			defer group.Done()
			t.Run(seed.ctx)
		}(seed.thread[i], seed.wg)
	}
}

// Wait ...
func (seed *seed) Wait() {
	seed.wg.Wait()
}

func defaultSeed() *seed {
	return &seed{
		Unfinished: make(map[string]*model.Unfinished),
		Videos:     make(map[string]*model.Video),
		Moves:      make(map[string]string),
		MaxLimit:   math.MaxUint16,
		wg:         &sync.WaitGroup{},
		thread:     make(map[Stepper]Threader, StepperMax),
		ignores:    make(map[string][]byte),
	}
}

// NewSeed ...
func NewSeed(ops ...Optioner) Seeder {
	seed := defaultSeed()
	seed.ctx, seed.cancel = context.WithCancel(context.Background())
	seed.Register(ops...)
	return seed
}

// Register ...
func (seed *seed) Register(ops ...Optioner) {
	for _, op := range ops {
		op.Option(seed)
	}
}

// Extend ...
type Extend struct {
	Path    string `json:"path"`
	Message string `json:"message"`
}

// Hash ...
func Hash(v interface{}) string {
	bytes, e := jsoniter.Marshal(v)
	if e != nil {
		return ""
	}
	return fmt.Sprintf("%x", sha1.Sum([]byte(bytes)))
}

// SkipTypeVerify ...
func SkipTypeVerify(tp string, v ...interface{}) bool {
	for i := range v {
		if v1, b := (v[i]).(string); b {
			if v1 == tp {
				return true
			}
		}
	}
	return false
}
