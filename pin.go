package seed

import (
	"context"
	shell "github.com/godcong/go-ipfs-restapi"
	"sync"

	"github.com/glvd/seed/model"
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
	skipSource bool
	state      PinStatus
	flag       PinFlag
	status     PinStatus
	list       []string
	index      int
	random     bool
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

	p.skipSource = seed.skipSource

}

// AfterRun ...
func (p *pin) AfterRun(seed *Seed) {
	return
}

// PinStatus ...
type PinStatus string

// PinStatusAll ...
const PinStatusAll PinStatus = "all"

// PinStatusUnfinished ...
const PinStatusUnfinished PinStatus = "unfinished"

// PinStatusAssignHash ...
const PinStatusAssignHash PinStatus = "assignHash"

// PinStatusAssignRelate ...
const PinStatusAssignRelate PinStatus = "assignRelate"

// PinStatusSliceOnly ...
const PinStatusSliceOnly PinStatus = "sliceOnly"

// PinStatusVideo ...
const PinStatusVideo PinStatus = "video"

// PinStatusVideo ...
const PinStatusPoster PinStatus = "poster"

// Pin ...
func Pin(status PinStatus, list ...string) Options {
	pin := &pin{
		status: status,
		list:   list,
		wg:     &sync.WaitGroup{},
	}

	return pinOption(pin)
}

// Run ...
func (p *pin) Run(ctx context.Context) {
	log.Info("pin running")
	switch p.status {
	case PinStatusAll:
		unfins, e := model.AllUnfinished(nil, 0)
		if e != nil {
			log.Error(e)
			return
		}

		log.Infof("pin(%d)", len(*unfins))
		for _, unf := range *unfins {
			select {
			case <-ctx.Done():
				return
			default:
				if p.skipSource && unf.Type == model.TypeVideo {
					log.With("type", unf.Type, "hash", unf.Hash, "sharpness", unf.Sharpness, "relate", unf.Relate).Info("pin skip")
					continue
				}

				log.With("type", unf.Type, "hash", unf.Hash, "sharpness", unf.Sharpness, "relate", unf.Relate).Info("pin")
				p.wg.Add(1)
				go p.pinHash(unf.Hash)
				p.wg.Wait()
				unf.Sync = true
				p.unfinished[unf.Hash] = unf
				e := model.AddOrUpdateUnfinished(unf)
				if e != nil {
					log.Error(e)
					continue
				}
			}
		}
	case PinStatusUnfinished:
		for hash, unf := range p.unfinished {
			select {
			case <-ctx.Done():
				return
			default:
				log.With("type", unf.Type, "hash", unf.Hash, "sharpness", unf.Sharpness, "relate", unf.Relate).Info("pin")
				p.wg.Add(1)
				go p.pinHash(hash)
				p.wg.Wait()
				p.unfinished[hash].Sync = true
				e := model.AddOrUpdateUnfinished(p.unfinished[hash])
				if e != nil {
					continue
				}
			}
		}
	case PinStatusAssignHash:
		for _, hash := range p.list {
			select {
			case <-ctx.Done():
				return
			default:
				p.wg.Add(1)
				go p.pinHash(hash)
				p.wg.Wait()
			}
		}
	case PinStatusAssignRelate:
		for _, relate := range p.list {
			select {
			case <-ctx.Done():
				return
			default:
				unfins, e := model.AllUnfinished(model.DB().Where("relate like ?", relate+"%"), 0)
				if e != nil {
					log.Error(e)
					continue
				}
				for _, unfin := range *unfins {
					p.wg.Add(1)
					go p.pinHash(unfin.Hash)
					p.wg.Wait()
				}
			}
		}
	case PinStatusSliceOnly:
		unfins, e := model.AllUnfinished(model.DB().Where("type = ?", model.TypeSlice), 0)
		if e != nil {
			log.Error(e)
			return
		}
		for _, unfin := range *unfins {
			select {
			case <-ctx.Done():
				return
			default:
				p.wg.Add(1)
				go p.pinHash(unfin.Hash)
				p.wg.Wait()
			}
		}
	case PinStatusVideo:
		videos, e := model.AllVideos(nil, 0)
		if e != nil {
			log.Error(e)
			return
		}
		for _, v := range *videos {
			log.With("bangumi", v.Bangumi, "poster", v.PosterHash, "m3u8", v.M3U8Hash, "thumb", v.ThumbHash, "source", v.SourceHash).Info("pin")
			if v.PosterHash != "" {
				p.wg.Add(1)
				go p.pinHash(v.PosterHash)
				p.wg.Wait()
			}

			if v.ThumbHash != "" {
				p.wg.Add(1)
				go p.pinHash(v.ThumbHash)
				p.wg.Wait()
			}

			if !p.skipSource && v.SourceHash != "" {
				p.wg.Add(1)
				go p.pinHash(v.SourceHash)
				p.wg.Wait()
			}
			if v.M3U8Hash != "" {
				p.wg.Add(1)
				go p.pinHash(v.M3U8Hash)
				p.wg.Wait()
			}

		}
	case PinStatusPoster:
		i, e := model.DB().Where("type = ?", model.TypePoster).Count(model.Unfinished{})
		if e != nil {
			log.Error(e)
			return
		}
		for start := 0; start < int(i); i += 50 {
			unfins, e := model.AllUnfinished(model.DB().Where("type = ?", model.TypePoster), 50, start)
			if e != nil {
				log.Error(e)
				return
			}

			log.Infof("pin(%d)", len(*unfins))
			for _, unf := range *unfins {
				select {
				case <-ctx.Done():
					return
				default:
					log.With("type", unf.Type, "hash", unf.Hash, "sharpness", unf.Sharpness, "relate", unf.Relate).Info("pin")
					p.wg.Add(1)
					go p.pinHash(unf.Hash)
					p.wg.Wait()
					unf.Sync = true
					p.unfinished[unf.Hash] = unf
					e := model.AddOrUpdateUnfinished(unf)
					if e != nil {
						log.Error(e)
						continue
					}
				}
			}
		}
	}
}

func (p *pin) pinHash(hash string) {
	defer func() {
		if p.wg != nil {
			p.wg.Done()
		}
	}()
	log.Info("pinning:", hash)
	e := p.shell.Pin(hash)
	if e != nil {
		log.Error("pin error:", hash, e)
		return
	}

}
