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
	thread      map[Stepper]Threader
	base        map[Stepper][]byte
	ignores     map[string][]byte
	err         error
	skipExist   bool
}

// GetThread ...
func (s *seed) GetThread(stepper Stepper) Threader {
	return s.thread[stepper]
}

// SetThread ...
func (s *seed) SetThread(stepper Stepper, threader Threader) {
	s.thread[stepper] = threader
}

// SetBaseThread ...
func (s *seed) SetBaseThread(stepper Stepper, threader Threader) {
	s.base[stepper] = nil
	s.thread[stepper] = threader
}

// IsBase ...
func (s *seed) IsBase(stepper Stepper) bool {
	_, b := s.base[stepper]
	return b
}

// HasThread ...
func (s *seed) HasThread(stepper Stepper) bool {
	_, b := s.thread[stepper]
	return b
}

// PushTo ...
func (s *seed) PushTo(stepper Stepper, v interface{}) (e error) {
	if val, b := s.thread[stepper]; b {
		return val.Push(v)
	}
	return xerrors.Errorf("thread(%d) is not exist", stepper)
}

// Args ...
func (s *seed) Args() map[string]interface{} {
	return s.args
}

// SetArgs ...
func (s *seed) SetArgs(args map[string]interface{}) {
	s.args = args
}

// AddArg ...
func (s *seed) AddArg(key string, value interface{}) {
	s.args[key] = value
}

// GetArg ...
func (s *seed) GetArg(key string) (v interface{}, b bool) {
	v, b = s.args[key]
	return
}

// GetStringArg ...
func (s *seed) GetStringArg(key string) (v string) {
	if arg, b := s.GetArg(key); b {
		v, _ = arg.(string)
	}
	return
}

// GetBoolArg ...
func (s *seed) GetBoolArg(key string) (v bool) {
	if arg, b := s.GetArg(key); b {
		v, _ = arg.(bool)
	}
	return
}

// GetNumberArg ...
func (s *seed) GetNumberArg(key string) (v int64) {
	if arg, b := s.GetArg(key); b {
		v, _ = arg.(int64)
	}
	return
}

// Done ...
func (s *seed) Done() {

}

// Stop ...
func (s *seed) Stop() {
	if s.cancel != nil {
		s.cancel()
	}
}

// Err ...
func (s *seed) Err() error {
	return s.err
}

// Start ...
func (s *seed) Start() {
	log.Info("Seed starting")
	for i := range s.thread {
		if s.thread[i] == nil {
			continue
		}
		s.thread[i].BeforeRun(s)

		s.wg.Add(1)
		go func(t Threader, s *seed) {
			defer s.wg.Done()
			t.Run(s.ctx)
			t.AfterRun(s)
		}(s.thread[i], s)
	}
}

// Wait ...
func (s *seed) Wait() {
	s.wg.Wait()

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
func (s *seed) Register(ops ...Optioner) {
	for _, op := range ops {
		op.Option(s)
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
