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
	unfinished []*model.Unfinished
	shell      *shell.Shell
	state      PinState
	flag       PinFlag
}

// BeforeRun ...
func (p *pin) BeforeRun(seed *Seed) {
	p.unfinished = seed.Unfinished
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
func Pin(flag PinFlag, shell ...*shell.Shell) Options {
	pin := &pin{
		flag: flag,
	}
	if shell != nil {
		pin.shell = shell[0]
	}

	return PinOption(pin)
}

// Run ...
func (p *pin) Run(ctx context.Context) {
	log.Infof("%+v", p.unfinished)
	wg := &sync.WaitGroup{}
	for _, v := range p.unfinished {
		select {
		case <-ctx.Done():
			return
		default:
			switch p.flag {
			case PinFlagSource:
				pinHash(nil, v.Hash)
			case PinFlagSlice:
				pinHash(nil, v.SliceHash)
			case PinFlagAll:
				if v.Hash != "" {
					wg.Add(1)
					go pinHash(wg, v.Hash)
				}
				if v.SliceHash != "" {
					wg.Add(1)
					go pinHash(wg, v.SliceHash)
				}
				wg.Wait()
			default:
				//nothing to do
			}
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
