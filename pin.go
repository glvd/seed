package seed

import (
	"context"
	shell "github.com/godcong/go-ipfs-restapi"
	"sync"

	"github.com/yinhevr/seed/model"
)

// PinFlag ...
type PinFlag string

// PinFlagNone ...
const (
	PinFlagNone PinFlag = "none"
	//PinFlagPoster PinFlag = "poster"
	PinFlagSource PinFlag = "source"
	PinFlagSlice  PinFlag = "slice"
	PinFlagAll    PinFlag = "all"
)

type pin struct {
	wg         *sync.WaitGroup
	unfinished map[string]*model.Unfinished
	shell      *shell.Shell
	state      PinState
	flag       PinFlag
}

// BeforeRun ...
func (p *pin) BeforeRun(seed *Seed) {
	p.unfinished = seed.Unfinished
	if p.unfinished == nil {
		p.unfinished = make(map[string]*model.Unfinished)
	}
	if p.shell == nil {
		p.shell = seed.Shell
	}

}

// AfterRun ...
func (p *pin) AfterRun(seed *Seed) {
	return
}

// PinState ...
type PinState string

// PinStateLocal ...
const PinStateLocal PinState = "local"

// PinStateRemote ...
const PinStateRemote PinState = "remote"

// Pin ...
func Pin(flag PinFlag) Options {
	pin := &pin{
		flag: flag,
		wg:   &sync.WaitGroup{},
	}

	return PinOption(pin)
}

// Run ...
func (p *pin) Run(ctx context.Context) {
	log.Infof("%+v", p.unfinished)
	for hash, unfin := range p.unfinished {
		log.Infof("%+v", unfin)
		select {
		case <-ctx.Done():
			return
		default:
			p.wg.Add(1)
			switch p.flag {
			case PinFlagSource:
				go p.pinHash(hash)
			case PinFlagSlice:
				go p.pinHash(hash)
			case PinFlagAll:
				go p.pinHash(hash)
			}
			p.wg.Wait()
		}
	}
}

func (p *pin) pinHash(hash string) {
	log.Info("pin:", hash)
	defer func() {
		if p.wg != nil {
			p.wg.Done()
		}
	}()
	e := p.shell.Pin(hash)
	if e != nil {
		log.Error("pin error:", hash, e)
		return
	}

	log.Info("pinned:", hash)
}
