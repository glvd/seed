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
)

// Options ...
type Options func(*Seed)

// AfterInitOptions ...
type AfterInitOptions func(*Seed)

// Thread ...
type Thread struct {
	wg sync.WaitGroup
}

// Threader ...
type Threader interface {
	Runnable
	BeforeRun(seed *Seed)
	AfterRun(seed *Seed)
}

// Runnable ...
type Runnable interface {
	Run(context.Context)
}

// Stepper ...
type Stepper int

// StepperNone ...
const (
	// StepperNone ...
	StepperNone Stepper = iota
	//StepperDatabase ...
	StepperDatabase
	// StepperInformation ...
	StepperInformation
	// StepperMoveInfo ...
	StepperMoveInfo
	// StepperProcess ...
	StepperProcess
	// StepperMoveproc ...
	StepperMoveproc
	// StepperTransfer ...
	StepperTransfer
	// StepperPin ...
	StepperPin
	// StepperCheck ...
	StepperCheck
	// StepperUpdate ...
	StepperUpdate
	// StepperMax ...
	StepperMax
)

// Seed ...
type Seed struct {
	Shell       *shell.Shell
	API         *API
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
	thread      []Threader
	ignores     map[string][]byte
	err         error
	skipExist   bool
}

// Args ...
func (seed *Seed) Args() map[string]interface{} {
	return seed.args
}

// SetArgs ...
func (seed *Seed) SetArgs(args map[string]interface{}) {
	seed.args = args
}

// AddArg ...
func (seed *Seed) AddArg(key string, value interface{}) {
	seed.args[key] = value
}

// GetArg ...
func (seed *Seed) GetArg(key string) interface{} {
	return seed.args[key]
}

// Stop ...
func (seed *Seed) Stop() {
	if seed.cancel != nil {
		seed.cancel()
	}
}

// Err ...
func (seed *Seed) Err() error {
	return seed.err
}

// Start ...
func (seed *Seed) Start() {
	go func() {
		log.Info("first running")
		defer seed.wg.Done()
		for i := range seed.thread {
			if seed.thread[i] == nil {
				continue
			}
			go func(t Threader, group *sync.WaitGroup) {
				defer group.Done()
				t.Run(seed.ctx)
			}(seed.thread[i], seed.wg)
		}
	}()
}

// Wait ...
func (seed *Seed) Wait() {
	seed.wg.Wait()
}

func defaultSeed() *Seed {
	return &Seed{
		Unfinished: make(map[string]*model.Unfinished),
		Videos:     make(map[string]*model.Video),
		Moves:      make(map[string]string),
		MaxLimit:   math.MaxUint16,
		wg:         &sync.WaitGroup{},
		threads:    3,
		thread:     make([]Threader, StepperMax),
		ignores:    make(map[string][]byte),
	}
}

// NewSeed ...
func NewSeed(ops ...Optioner) *Seed {
	seed := defaultSeed()
	seed.ctx, seed.cancel = context.WithCancel(context.Background())
	seed.Register(ops...)
	return seed
}

// Register ...
func (seed *Seed) Register(ops ...Optioner) {
	for _, op := range ops {
		op.Option(seed)
	}
}

// SkipConvertOption ...
func SkipConvertOption() Options {
	return func(seed *Seed) {
		seed.skipConvert = true
	}
}

// SkipExistOption ...
func SkipExistOption() Options {
	return func(seed *Seed) {
		seed.skipExist = true
	}
}

// NoSliceOption ...
func NoSliceOption() Options {
	return func(seed *Seed) {
		seed.noSlice = true
	}
}

// MaxLimitOption ...
func MaxLimitOption(max int) Options {
	return func(seed *Seed) {
		seed.MaxLimit = max
	}
}

// PreAddOption ...
func PreAddOption() Options {
	return func(seed *Seed) {
		seed.preAdd = true
	}
}

// databaseOption ...
func databaseOption(db *Database) Options {
	return func(seed *Seed) {
		seed.thread[StepperDatabase] = db
	}
}

// InformationOption ...
func informationOption(info *Information) Options {
	return func(seed *Seed) {
		seed.thread[StepperInformation] = info
	}
}

// UnfinishedOption ...
func UnfinishedOption(unfins ...*model.Unfinished) Options {
	return func(seed *Seed) {
		if seed.Unfinished == nil {
			seed.Unfinished = make(map[string]*model.Unfinished)
		}
		for _, u := range unfins {
			if u == nil {
				continue
			}

			if u.Hash != "" {
				seed.Unfinished[u.Hash] = u
			}
		}
	}
}

// ShellOption ...
func ShellOption(s string) Options {
	return func(seed *Seed) {
		log.Info("ipfs: ", s)
		seed.Shell = shell.NewShell(s)
	}
}

// APIOption ...
func APIOption(s string) Options {
	return func(seed *Seed) {
		seed.API = NewAPI(s)
	}
}

// processOption ...
func processOption(process *process) Options {
	return func(seed *Seed) {
		seed.thread[StepperProcess] = process
	}
}

// pinOption ...
func pinOption(pin *pin) Options {
	return func(seed *Seed) {
		seed.thread[StepperPin] = pin
	}
}

// IgnoreOption ...
func IgnoreOption(ignores ...string) Options {
	return func(seed *Seed) {
		for _, i := range ignores {
			seed.ignores[i] = nil
		}
	}
}

// ThreadOption ...
func ThreadOption(t int) Options {
	return func(seed *Seed) {
		seed.threads = t
	}
}

// Extend ...
type Extend struct {
	Path    string `json:"path"`
	Message string `json:"message"`
}

// VideoSource ...

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
