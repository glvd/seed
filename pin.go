package seed

import (
	"context"
	"sync"

	httpapi "github.com/ipfs/go-ipfs-http-client"

	shell "github.com/godcong/go-ipfs-restapi"

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
	api        *httpapi.HttpApi
	wg         *sync.WaitGroup
	unfinished map[string]*model.Unfinished
	shell      *shell.Shell
	skipType   []interface{}
	skipSource bool
	state      PinStatus
	flag       PinFlag
	status     PinStatus
	list       []string
	index      int
	random     bool
	from       string
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
	p.from = seed.From
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

// PinStatusVideo ...
const PinStatusVideo PinStatus = "video"

// PinStatusPoster ...
const PinStatusPoster PinStatus = "poster"

// PinStatusSync ...
const PinStatusSync PinStatus = "sync"

// PinArgs ...
type PinArgs func(c *pin)

// PinSkipArg ...
func PinSkipArg(s []string) PinArgs {
	return func(p *pin) {
		if s == nil {
			return
		}
		for i := range s {
			p.skipType = append(p.skipType, s[i])
		}
	}
}

// PinListArg ...
func PinListArg(s ...string) PinArgs {
	return func(p *pin) {
		p.list = s
	}
}

// PinStatusArg ...
func PinStatusArg(s PinStatus) PinArgs {
	return func(p *pin) {
		p.status = s
	}
}

// Pin ...
func Pin(args ...PinArgs) Options {
	pin := &pin{
		status: PinStatusAll,
	}

	for _, argFn := range args {
		argFn(pin)
	}
	return pinOption(pin)
}

// Run ...
func (p *pin) Run(ctx context.Context) {
	log.Info("pin running")
	switch p.status {
	case PinStatusAll:
		s := model.DB().NewSession()
		if len(p.skipType) > 0 {
			s.NotIn("type", p.skipType...)
		}
		i, e := s.Clone().Count(model.Unfinished{})
		if e != nil {
			log.Error(e)
			return
		}
		for start := 0; start < int(i); start += 50 {
			unfins, e := model.AllUnfinished(s.Clone(), 50, start)
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
					e := p.pinHash(unf.Hash)
					if e != nil {
						log.Error(e)
						return
					}
					unf.Sync = true
					p.unfinished[unf.Hash] = unf
					e = model.AddOrUpdateUnfinished(unf)
					if e != nil {
						log.Error(e)
						continue
					}
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
				e := p.pinHash(hash)
				if e != nil {
					log.Error(e)
					return
				}
				p.unfinished[hash].Sync = true
				e = model.AddOrUpdateUnfinished(p.unfinished[hash])
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
				e := p.pinHash(hash)
				if e != nil {
					log.Error(e)
					return
				}
			}
		}
	case PinStatusAssignRelate:
		for _, relate := range p.list {
			select {
			case <-ctx.Done():
				return
			default:
				unfins, e := model.AllUnfinished(model.DB().Where("relate = ?", relate).Or("relate like ?", relate+"-%"), 0)
				if e != nil {
					log.Error(e)
					continue
				}
				for _, unfin := range *unfins {
					e := p.pinHash(unfin.Hash)
					if e != nil {
						log.Error(e)
						return
					}
				}
			}
		}
	case PinStatusVideo:
		i, e := model.DB().Count(model.Video{})
		if e != nil {
			log.Error(e)
			return
		}
		for start := 0; start < int(i); start += 50 {
			videos, e := model.AllVideos(nil, 50, start)
			if e != nil {
				log.Error(e)
				return
			}
			for _, v := range *videos {
				log.With("bangumi", v.Bangumi, "poster", v.PosterHash, "m3u8", v.M3U8Hash, "thumb", v.ThumbHash, "source", v.SourceHash).Info("pin")

				if !SkipTypeVerify("poster", p.skipType...) && v.PosterHash != "" {
					e := p.pinHash(v.PosterHash)
					if e != nil {
						log.Error(e)
						return
					}
				}

				if !SkipTypeVerify("thumb", p.skipType...) && v.ThumbHash != "" {
					e := p.pinHash(v.ThumbHash)
					if e != nil {
						log.Error(e)
						return
					}
				}

				if !SkipTypeVerify("source", p.skipType...) && v.SourceHash != "" {
					e := p.pinHash(v.SourceHash)
					if e != nil {
						log.Error(e)
						return
					}
				}
				if !SkipTypeVerify("slice", p.skipType...) && v.M3U8Hash != "" {
					e := p.pinHash(v.M3U8Hash)
					if e != nil {
						log.Error(e)
						return
					}
				}

			}
		}

	case PinStatusSync:
		s := model.DB().Where("machine_id like ?", "%"+p.from+"%")
		i, e := s.Clone().Count(model.Pin{})
		if e != nil {
			log.Error(e)
			return
		}
		for start := 0; start < int(i); start += 50 {
			pins, e := model.AllPin(s.Clone(), 50, start)
			if e != nil {
				log.Error(e)
				return
			}
			for _, ps := range *pins {
				select {
				case <-ctx.Done():
					return
				default:
					e := p.pinHash(ps.PinHash)
					if e != nil {
						log.Error(e)
						return
					}
				}
			}
		}
	}
}

func (p *pin) pinHash(hash string) (e error) {
	log.Info("pinning:", hash)
	e = p.shell.Pin(hash)
	if e != nil {
		log.With("hash", hash).Error(e)
	}
	return
}
