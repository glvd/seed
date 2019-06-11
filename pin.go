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
	}

	return PinOption(pin)
}

// Run ...
func (p *pin) Run(ctx context.Context) {
	log.Infof("%+v", p.unfinished)
	wg := &sync.WaitGroup{}
	for hash, unfin := range p.unfinished {
		log.Infof("%+v", unfin)
		select {
		case <-ctx.Done():
			return
		default:
			wg.Add(1)
			switch p.flag {
			case PinFlagSource:
				go pinHash(wg, hash)
			case PinFlagSlice:
				go pinHash(wg, hash)
			case PinFlagAll:
				go pinHash(wg, hash)
			}
			wg.Wait()
		}
	}
}

func pinHash(wg *sync.WaitGroup, hash string) {
	log.Info("pin:", hash)
	defer func() {
		if wg != nil {
			wg.Done()
		}
	}()
	e := rest.Pin(hash)
	if e != nil {
		log.Error("pin error:", hash, e)
		return
	}

	log.Info("pinned:", hash)
}
