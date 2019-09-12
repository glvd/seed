package task

import (
	"context"

	"github.com/glvd/seed"
	"github.com/glvd/seed/model"
	"github.com/go-xorm/xorm"
	httpapi "github.com/ipfs/go-ipfs-http-client"
	iface "github.com/ipfs/interface-go-ipfs-core"
	"github.com/ipfs/interface-go-ipfs-core/options"
	"github.com/ipfs/interface-go-ipfs-core/path"
)

// Pin ...
type Pin struct {
	Table PinTable
	//CheckType  CheckType
	SkipType []interface{}
	Type     PinType
	list     []string
	index    int
	random   bool
	from     string
}

// CallTask ...
func (p *Pin) CallTask(seeder seed.Seeder, task *seed.Task) error {
	select {
	case <-seeder.Context().Done():
		return nil
	default:
		switch p.Type {
		case PinTypeAdd:
			pin := &pinAdd{table: p.Table}
			e := seeder.PushTo(seed.StepperAPI, pin)
			if e != nil {
				log.Error(e)
				return e
			}
		}

	}
	return nil
}

// Task ...
func (p *Pin) Task() *seed.Task {
	return seed.NewTask(p)
}

// PinTable ...
type PinTable string

// PinTableUnfinished ...
const PinTableUnfinished PinTable = "unfinished"

// PinTableVideo ...
const PinTableVideo PinTable = "video"

// PinType ...
type PinType string

// PinTypeCheck ...
const PinTypeCheck PinType = "check"

// PinTypeAdd ...
const PinTypeAdd PinType = "add"

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

// NewPin ...
func NewPin(args ...PinArgs) *Pin {
	pin := &Pin{
		Table: PinTableUnfinished,
		Type:  PinTypeAdd,
	}
	for _, argFn := range args {
		argFn(pin)
	}
	return pin
}

type pinAdd struct {
	table PinTable
	skip  []interface{}
}

func (p *pinAdd) pinUnfinishedCall(a *seed.API, api *httpapi.HttpApi) {
	u := make(chan *model.Unfinished)
	e := a.PushTo(seed.UnfinishedCall(u, func(session *xorm.Session) *xorm.Session {
		return session
	}))
	if e != nil {
		log.Error(e)
	}
ChanEnd:
	for {
		select {
		case unfinished := <-u:
			if unfinished == nil {
				break ChanEnd
			}
			if !seed.SkipTypeVerify(unfinished.Type, p.skip...) {
				log.With("type", unfinished.Type, "hash", unfinished.Hash).Info("pinning")
				e := api.Pin().Add(a.Context(), path.New(unfinished.Hash), func(settings *options.PinAddSettings) error {
					settings.Recursive = true
					return nil
				})
				if e != nil {
					log.Error(e)
					break ChanEnd
				}
			}
		}
	}
	close(u)
}

func (p *pinAdd) pinVideoCall(a *seed.API, api *httpapi.HttpApi) {
	v := make(chan *model.Video)
	e := a.PushTo(seed.VideoCall(v, func(session *xorm.Session) *xorm.Session {
		return session
	}))
	if e != nil {
		log.Error(e)
	}
ChanEnd:
	for {
		select {
		case video := <-v:
			if video == nil {
				break ChanEnd
			}
			if !seed.SkipVerify("source", p.skip...) && video.SourceHash != "" {
				log.With("hash", video.SourceHash).Info("source pinning")
				e := api.Pin().Add(a.Context(), path.New(video.SourceHash), func(settings *options.PinAddSettings) error {
					settings.Recursive = true
					return nil
				})
				if e != nil {
					log.Error(e)
					break ChanEnd
				}
			}
			if !seed.SkipVerify("slice", p.skip...) && video.M3U8Hash != "" {
				log.With("hash", video.M3U8Hash).Info("slice pinning")
				e := api.Pin().Add(a.Context(), path.New(video.M3U8Hash), func(settings *options.PinAddSettings) error {
					settings.Recursive = true
					return nil
				})
				if e != nil {
					log.Error(e)
					break ChanEnd
				}
			}
			if !seed.SkipVerify("poster", p.skip...) && video.PosterHash != "" {
				log.With("hash", video.PosterHash).Info("poster pinning")
				e := api.Pin().Add(a.Context(), path.New(video.PosterHash), func(settings *options.PinAddSettings) error {
					settings.Recursive = true
					return nil
				})
				if e != nil {
					log.Error(e)
					break ChanEnd
				}
			}
			if !seed.SkipVerify("thumb", p.skip...) && video.ThumbHash != "" {
				log.With("hash", video.ThumbHash).Info("thumb pinning")
				e := api.Pin().Add(a.Context(), path.New(video.ThumbHash), func(settings *options.PinAddSettings) error {
					settings.Recursive = true
					return nil
				})
				if e != nil {
					log.Error(e)
					break ChanEnd
				}
			}

		}
	}
	close(v)
}

// Call ...
func (p *pinAdd) Call(a *seed.API, api *httpapi.HttpApi) error {
	log.Info("pin add")
	if p.table == PinTableUnfinished {
		p.pinUnfinishedCall(a, api)
	} else if p.table == PinTableVideo {
		p.pinVideoCall(a, api)
	} else {
		//do nothing
	}

	return nil
}

type pinCheck struct {
	table     PinTable
	skip      []interface{}
	checkType CheckType
	checkOut  string
}

//Call ...
func (p *pinCheck) Call(a *seed.API, api *httpapi.HttpApi) error {
	log.Info("pin add")
	if p.table == PinTableUnfinished {
		p.pinUnfinishedCall(a, api)
	} else if p.table == PinTableVideo {
		p.pinVideoCall(a, api)
	} else {
		//do nothing
	}

	return nil
}

func (p *pinCheck) pinVideoCall(a *seed.API, api *httpapi.HttpApi) {

}

func (p *pinCheck) pinUnfinishedCall(a *seed.API, api *httpapi.HttpApi) {
	u := make(chan *model.Unfinished)
	e := a.PushTo(seed.UnfinishedCall(u, func(session *xorm.Session) *xorm.Session {
		return session
	}))
	if e != nil {
		log.Error(e)
	}
	pinned := make(map[string]*model.Unfinished)
	pins, e := api.Pin().Ls(a.Context(), func(settings *options.PinLsSettings) error {
		settings.Type = "recursive"
		return nil
	})

	if e != nil {
		log.Error(e)
	}
	for _, p := range pins {
		pinned[model.PinHash(p.Path())] = nil
	}

ChanEnd:
	for {
		select {
		case unfinished := <-u:
			if unfinished == nil {
				break ChanEnd
			}
			if !seed.SkipTypeVerify(unfinished.Type, p.skip...) {
				if _, b := pinned[unfinished.Hash]; b {
					if p.checkType == CheckTypePin || p.checkType == CheckTypeAll {
						log.With("hash", unfinished.Hash, "relate", unfinished.Relate, "type", unfinished.Type).Info("pinned")
					}
					pinned[unfinished.Hash] = unfinished
				} else {
					if p.checkType == CheckTypeUnpin || p.checkType == CheckTypeAll {
						log.With("hash", unfinished.Hash, "relate", unfinished.Relate, "type", unfinished.Type).Info("unpin")
					}
				}

			}
		}
	}
	close(u)
}

func listPin(ctx context.Context, p *Pin) <-chan iface.Pin {
	cPin := make(chan iface.Pin)
	//(seed.StepperAPI, func(api *seed.API, api2 *httpapi.HttpApi) (e error) {
	//	defer func() {
	//		cPin <- nil
	//	}()
	//	pins, e := api2.Pin().Ls(ctx, func(settings *options.PinLsSettings) error {
	//		settings.Type = p.Type
	//		return nil
	//	})
	//	if e != nil {
	//		return e
	//	}
	//	for _, p := range pins {
	//		cPin <- p
	//	}
	//
	//	return nil
	//})
	//if e != nil {
	//	return nil
	//}
	return cPin
}

func unfinishedList(ctx context.Context, p *Pin) <-chan *model.Unfinished {
	u := make(chan *model.Unfinished)
	//defer func() {
	//	u <- nil
	//}()
	//p.PushTo(seed.StepperDatabase, func(database *seed.Database, eng *xorm.Engine) (e error) {
	//	session := eng.NewSession()
	//	if len(p.SkipType) > 0 {
	//		session = session.NotIn("type", p.SkipType...)
	//	}
	//	i, e := session.Clone().Count(model.Unfinished{})
	//	if e != nil {
	//		return e
	//	}
	//	for start := 0; start < int(i); start += 50 {
	//		unfinishs, e := model.AllUnfinished(session.Clone(), 50, start)
	//		if e != nil {
	//			return e
	//		}
	//		for i := range *unfinishs {
	//			u <- (*unfinishs)[i]
	//		}
	//	}
	//	return nil
	//})
	//
	return u
}

// Run ...
func (p *Pin) Run(ctx context.Context) {
	log.Info("pin running")
	//switch p.PinType {
	//case PinTypeAdd:
	//	//pin := listPin(ctx, p.API, p.Type)
	//case PinTypeCheck:
	//	//pin := listPin(ctx, p.API, p.Type)
	//}
	//
	//switch p.PinStatus {
	//case PinStatusAll:
	//	u := unfinishedList(ctx, p)
	//	for {
	//		select {
	//		case <-ctx.Done():
	//			return
	//		case uf := <-u:
	//			log.With("type", uf.Type, "hash", uf.Hash, "sharpness", uf.Sharpness, "relate", uf.Relate).Info("Pin")
	//			//if e := APIPin(p.Seeder, uf.Hash); e != nil {
	//			//	log.With("hash", uf.Hash).Error(e, " not pinned")
	//			//}
	//		}
	//	}
	//case PinStatusAssignHash:
	//	for _, hash := range p.list {
	//		select {
	//		case <-ctx.Done():
	//			return
	//		default:
	//			e := p.pinHash(hash)
	//			if e != nil {
	//				log.Error(e)
	//				return
	//			}
	//		}
	//	}
	//case PinStatusAssignRelate:
	//	for _, relate := range p.list {
	//		select {
	//		case <-ctx.Done():
	//			return
	//		default:
	//			p.Database.PushCallback(func(database *Database, eng *xorm.Engine) (e error) {
	//				unfins, e := model.AllUnfinished(eng.Where("relate = ?", relate).Or("relate like ?", relate+"-%"), 0)
	//				if e != nil {
	//					return e
	//				}
	//				for _, unfin := range *unfins {
	//					e := p.pinHash(unfin.Hash)
	//					if e != nil {
	//						return e
	//					}
	//				}
	//				return nil
	//			})
	//
	//		}
	//	}
	//case PinTableVideo:
	//	i, e := model.DB().Count(model.Video{})
	//	if e != nil {
	//		log.Error(e)
	//		return
	//	}
	//	for start := 0; start < int(i); start += 50 {
	//		videos, e := model.AllVideos(nil, 50, start)
	//		if e != nil {
	//			log.Error(e)
	//			return
	//		}
	//		for _, v := range *videos {
	//			log.With("bangumi", v.Bangumi, "poster", v.PosterHash, "m3u8", v.M3U8Hash, "thumb", v.ThumbHash, "source", v.SourceHash).Info("Pin")
	//
	//			if !SkipVerify("poster", p.SkipType...) && v.PosterHash != "" {
	//				e := p.pinHash(v.PosterHash)
	//				if e != nil {
	//					log.Error(e)
	//					return
	//				}
	//			}
	//
	//			if !SkipVerify("thumb", p.SkipType...) && v.ThumbHash != "" {
	//				e := p.pinHash(v.ThumbHash)
	//				if e != nil {
	//					log.Error(e)
	//					return
	//				}
	//			}
	//
	//			if !SkipVerify("source", p.SkipType...) && v.SourceHash != "" {
	//				e := p.pinHash(v.SourceHash)
	//				if e != nil {
	//					log.Error(e)
	//					return
	//				}
	//			}
	//			if !SkipVerify("slice", p.SkipType...) && v.M3U8Hash != "" {
	//				e := p.pinHash(v.M3U8Hash)
	//				if e != nil {
	//					log.Error(e)
	//					return
	//				}
	//			}
	//
	//		}
	//	}
	//
	//case PinStatusSync:
	//	s := model.DB().Where("machine_id like ?", "%"+p.from+"%")
	//	i, e := s.Clone().Count(model.Pin{})
	//	if e != nil {
	//		log.Error(e)
	//		return
	//	}
	//	for start := 0; start < int(i); start += 50 {
	//		pins, e := model.AllPin(s.Clone(), 50, start)
	//		if e != nil {
	//			log.Error(e)
	//			return
	//		}
	//		for _, ps := range *pins {
	//			select {
	//			case <-ctx.Done():
	//				return
	//			default:
	//				e := p.pinHash(ps.PinHash)
	//				if e != nil {
	//					log.Error(e)
	//					return
	//				}
	//			}
	//		}
	//	}
	//}
	//}
}

//func (p *Pin) pinHash(hash string) (e error) {
//	log.Info("pinning:", hash)
//	e = p.Pin(hash)
//	if e != nil {
//		log.With("hash", hash).Error(e)
//	}
//	return
//}

// PinCallFunc ...
type PinCallFunc func(*Pin) error

// PinCaller ...
type PinCaller interface {
	Call()
}

type pinCall struct {
}
