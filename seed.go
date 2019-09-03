package seed

import (
	"context"
	"crypto/sha1"
	"fmt"
	"sync"
	"time"

	json "github.com/json-iterator/go"
	"go.uber.org/atomic"
)

// Options ...
type Options func(seeder Seeder)

// AfterInitOptions ...
type AfterInitOptions func(Seeder)

// seed ...
type seed struct {
	args   map[string]interface{}
	wg     *sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
	thread map[Stepper]ThreadRun
	base   map[Stepper]ThreadBase
	normal map[Stepper][]byte
}

// AddTasker ...
func (s *seed) AddTasker(tasker Tasker) {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		e := tasker.Task().Push(s)
		if e != nil {
			log.Error(e)
		}
	}()
}

// Context ...
func (s *seed) Context() context.Context {
	return s.ctx
}

// RunTask ...
func (s *seed) RunTask(task *Task) {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		e := task.Push(s)
		if e != nil {
			log.Error(e)
		}
	}()
}

func defaultSeed() *seed {
	return &seed{
		wg:     &sync.WaitGroup{},
		thread: make(map[Stepper]ThreadRun, StepperMax),
		base:   make(map[Stepper]ThreadBase, StepperMax),
		normal: make(map[Stepper][]byte, StepperMax),
	}
}

// NewSeed ...
func NewSeed(ops ...Optioner) Seeder {
	seed := defaultSeed()
	seed.ctx, seed.cancel = context.WithCancel(context.Background())
	seed.Register(ops...)
	return seed
}

// GetThread ...
func (s *seed) GetThread(stepper Stepper) ThreadRun {
	return s.thread[stepper]
}

// SetThread ...
func (s *seed) SetThread(stepper Stepper, threader ThreadRun) {
	s.thread[stepper] = threader
}

// SetBaseThread ...
func (s *seed) SetBaseThread(stepper Stepper, threader Threader) {
	s.base[stepper] = threader
	s.thread[stepper] = threader
}

// IsBase ...
func (s *seed) IsBase(stepper Stepper) bool {
	_, b := s.base[stepper]
	return b
}

// SetNormalThread ...
func (s *seed) SetNormalThread(stepper Stepper, threader ThreadRun) {
	s.normal[stepper] = nil
	s.thread[stepper] = threader
}

// IsNormal ...
func (s *seed) IsNormal(stepper Stepper) bool {
	_, b := s.normal[stepper]
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
	return fmt.Errorf("thread(%d) is not exist", stepper)
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
	count := atomic.NewInt32(0)
	for i := range s.base {
		go func(base ThreadBase) {
			//if base.State() != StateStop {
			<-base.Done()
			//}
			count.Add(1)
		}(s.base[i])
		s.base[i].Finished()
	}
	for {
		if count.Load() == int32(len(s.base)) {
			return
		}
	}

}

// Stop ...
func (s *seed) Stop() {
	if s.cancel != nil {
		s.cancel()
	}
}

// Start ...
func (s *seed) Start() {
	log.Info("Seed starting")
	for i := range s.thread {
		if s.thread[i] == nil {
			continue
		}
		s.thread[i].BeforeRun(s)
		if s.IsNormal(i) || s.IsBase(i) {
			log.With("thread", i).Info("run base/normal")
			go func(t ThreadRun, s *seed) {
				t.Run(s.ctx)
				t.AfterRun(s)
			}(s.thread[i], s)
			continue
		} else {
			s.wg.Add(1)
			log.With("thread", i).Info("run thread")
			go func(t ThreadRun, s *seed) {
				defer s.wg.Done()
				t.Run(s.ctx)
				t.AfterRun(s)
			}(s.thread[i], s)
		}
	}
}

// Wait ...
func (s *seed) Wait() {
	s.wg.Wait()

	log.Info("waiting base")
	state := 1
	for state > 0 {
		for _, t := range s.base {
			if StateWaiting != t.State() {
				state = 2
				time.Sleep(15 * time.Second)
				break
			}
		}
		if state == 2 {
			state = 1
		} else {
			state = 0
		}
	}
	log.Info("base done")
	s.Done()
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
	bytes, e := json.Marshal(v)
	if e != nil {
		return ""
	}
	return fmt.Sprintf("%x", sha1.Sum(bytes))
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
