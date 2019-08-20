package seed

import (
	"context"

	"github.com/go-xorm/xorm"
	shell "github.com/godcong/go-ipfs-restapi"
	httpapi "github.com/ipfs/go-ipfs-http-client"
	iface "github.com/ipfs/interface-go-ipfs-core"
	"github.com/ipfs/interface-go-ipfs-core/options"

	"github.com/glvd/seed/model"
)

// PinFlag ...
type PinFlag string

// PinFlagNone ...
const (
	// PinFlagNone ...
	PinFlagNone PinFlag = "none"
	//PinFlagPoster PinFlag = "poster"
	PinFlagSource PinFlag = "source"
	// PinFlagSlice ...
	PinFlagSlice PinFlag = "slice"
	// PinFlagAll ...
	PinFlagAll PinFlag = "all"
)

// Pin ...
type Pin struct {
	*Seed
	PinType    PinType
	CheckType  CheckType
	SkipType   []interface{}
	Type       string
	PinStatus  PinStatus
	unfinished map[string]*model.Unfinished
	shell      *shell.Shell
	state      PinStatus
	flag       PinFlag
	list       []string
	index      int
	random     bool
	from       string
}

// BeforeRun ...
func (p *Pin) BeforeRun(seed *Seed) {
	p.Seed = seed
}

// AfterRun ...
func (p *Pin) AfterRun(seed *Seed) {
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

// PinType ...
type PinType string

// PinTypeCheck ...
const PinTypeCheck PinType = "check"

// PinTypeAdd ...
const PinTypeAdd PinType = "add"

// CheckType ...
type CheckType string

// CheckTypePin ...
const CheckTypePin CheckType = "Pin"

// CheckTypeUnpin ...
const CheckTypeUnpin CheckType = "unpin"

// PinArgs ...
type PinArgs func(c *Pin)

// PinSkipArg ...
func PinSkipArg(s []string) PinArgs {
	return func(p *Pin) {
		if s == nil {
			return
		}
		for i := range s {
			p.SkipType = append(p.SkipType, s[i])
		}
	}
}

// PinListArg ...
func PinListArg(s ...string) PinArgs {
	return func(p *Pin) {
		p.list = s
	}
}

// PinStatusArg ...
func PinStatusArg(s PinStatus) PinArgs {
	return func(p *Pin) {
		p.PinStatus = s
	}
}

// NewPin ...
func NewPin(args ...PinArgs) Options {
	pin := &Pin{
		PinStatus: PinStatusAll,
	}

	for _, argFn := range args {
		argFn(pin)
	}
	return pinOption(pin)
}

func listPin(ctx context.Context, api *API, t string) <-chan iface.Pin {
	cPin := make(chan iface.Pin)
	api.PushRun(func(api *API, api2 *httpapi.HttpApi) (e error) {
		defer func() {
			cPin <- nil
		}()
		pins, e := api2.Pin().Ls(ctx, func(settings *options.PinLsSettings) error {
			settings.Type = t
			return nil
		})
		if e != nil {
			return e
		}
		for _, p := range pins {
			cPin <- p
		}

		return nil
	})
	return cPin
}

func unfinishedList(ctx context.Context, engine *xorm.Engine) <-chan *model.Unfinished {

}

// Run ...
func (p *Pin) Run(ctx context.Context) {
	log.Info("Pin running")
	switch p.PinType {
	case PinTypeAdd:
		pin := listPin(ctx, p.API, p.Type)
	case PinTypeCheck:
		pin := listPin(ctx, p.API, p.Type)
	}

	switch p.PinStatus {
	case PinStatusAll:
		s := model.DB().NewSession()
		if len(p.SkipType) > 0 {
			s.NotIn("type", p.SkipType...)
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

			log.Infof("Pin(%d)", len(*unfins))
			for _, unf := range *unfins {
				select {
				case <-ctx.Done():
					return
				default:
					log.With("type", unf.Type, "hash", unf.Hash, "sharpness", unf.Sharpness, "relate", unf.Relate).Info("Pin")
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
				log.With("type", unf.Type, "hash", unf.Hash, "sharpness", unf.Sharpness, "relate", unf.Relate).Info("Pin")
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
				log.With("bangumi", v.Bangumi, "poster", v.PosterHash, "m3u8", v.M3U8Hash, "thumb", v.ThumbHash, "source", v.SourceHash).Info("Pin")

				if !SkipTypeVerify("poster", p.SkipType...) && v.PosterHash != "" {
					e := p.pinHash(v.PosterHash)
					if e != nil {
						log.Error(e)
						return
					}
				}

				if !SkipTypeVerify("thumb", p.SkipType...) && v.ThumbHash != "" {
					e := p.pinHash(v.ThumbHash)
					if e != nil {
						log.Error(e)
						return
					}
				}

				if !SkipTypeVerify("source", p.SkipType...) && v.SourceHash != "" {
					e := p.pinHash(v.SourceHash)
					if e != nil {
						log.Error(e)
						return
					}
				}
				if !SkipTypeVerify("slice", p.SkipType...) && v.M3U8Hash != "" {
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

func (p *Pin) pinHash(hash string) (e error) {
	log.Info("pinning:", hash)
	e = p.shell.Pin(hash)
	if e != nil {
		log.With("hash", hash).Error(e)
	}
	return
}
