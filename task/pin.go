package task

import (
	"context"
	"strings"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiformats/go-multiaddr"

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
	Table    PinTable
	Check    CheckType
	SkipType []interface{}
	Type     PinType
	list     []string
	index    int
	random   bool
	from     string
	From     string
}

// CallTask ...
func (p *Pin) CallTask(seeder seed.Seeder, task *seed.Task) error {
	select {
	case <-seeder.Context().Done():
		return nil
	default:
		switch p.Type {
		case PinTypeAdd:
			pin := &pinAdd{table: p.Table, skip: p.SkipType}
			e := seeder.PushTo(seed.StepperAPI, pin)
			if e != nil {
				log.Error(e)
				return e
			}
		case PinTypeCheck:
			pin := &pinCheck{table: p.Table, skip: p.SkipType, checkType: p.Check}
			e := seeder.PushTo(seed.StepperAPI, pin)
			if e != nil {
				log.Error(e)
				return e
			}
		case PinTypeSync:
			pin := &pinSync{table: p.Table, from: p.From, skip: p.SkipType}
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

// CheckType ...
type CheckType string

//CheckTypeAll ...
const CheckTypeAll CheckType = "all"

// CheckTypePin ...
const CheckTypePin CheckType = "pin"

// CheckTypeUnpin ...
const CheckTypeUnpin CheckType = "unpin"

// PinTable ...
type PinTable string

// PinTableUnfinished ...
const PinTableUnfinished PinTable = "unfinished"

// PinTableVideo ...
const PinTableVideo PinTable = "video"

//PinTablePin ...
const PinTablePin PinTable = "pin"

// PinType ...
type PinType string

// PinTypeCheck ...
const PinTypeCheck PinType = "check"

// PinTypeAdd ...
const PinTypeAdd PinType = "add"

//PinTypeSync ...
const PinTypeSync PinType = "sync"

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
		Check: CheckTypeAll,
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
	e := a.PushTo(seed.DatabaseUnfinishedCall(u, func(session *xorm.Session) *xorm.Session {
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
	e := a.PushTo(seed.DatabaseVideoCall(v, func(session *xorm.Session) *xorm.Session {
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
	log.Info("pin check")
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
	v := make(chan *model.Video)
	e := a.PushTo(seed.DatabaseVideoCall(v, func(session *xorm.Session) *xorm.Session {
		return session.Where("m3u8_hash <> ?", "")
	}))
	if e != nil {
		log.Error(e)
	}
	pinned := make(map[string]*model.Video)
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
	log.With("total", len(pinned)).Info("pinned")

	myid, e := seed.MyID(a)
	if e != nil {
		return
	}
ChanEnd:
	for {
		select {
		case video := <-v:
			if video == nil {
				break ChanEnd
			}
			if !seed.SkipTypeVerify(model.TypeThumb, p.skip...) {
				if _, b := pinned[video.ThumbHash]; b {
					if p.checkType == CheckTypePin || p.checkType == CheckTypeAll {
						log.With("hash", video.ThumbHash, "relate", video.Bangumi, "type", "thumb").Info("pinned")

					}
					pinned[video.ThumbHash] = video
				} else {
					if p.checkType == CheckTypeUnpin || p.checkType == CheckTypeAll {
						log.With("hash", video.ThumbHash, "relate", video.Bangumi, "type", "thumb").Info("unpin")
					}
				}

			}
			if !seed.SkipTypeVerify(model.TypePoster, p.skip...) {
				if _, b := pinned[video.PosterHash]; b {
					if p.checkType == CheckTypePin || p.checkType == CheckTypeAll {
						log.With("hash", video.PosterHash, "relate", video.Bangumi, "type", "poster").Info("pinned")

					}
					pinned[video.PosterHash] = video
				} else {
					if p.checkType == CheckTypeUnpin || p.checkType == CheckTypeAll {
						log.With("hash", video.PosterHash, "relate", video.Bangumi, "type", "poster").Info("unpin")
					}
				}

			}
			if !seed.SkipTypeVerify(model.TypeVideo, p.skip...) {
				if _, b := pinned[video.SourceHash]; b {
					if p.checkType == CheckTypePin || p.checkType == CheckTypeAll {
						log.With("hash", video.SourceHash, "relate", video.Bangumi, "type", "video").Info("pinned")

					}
					pinned[video.SourceHash] = video
				} else {
					if p.checkType == CheckTypeUnpin || p.checkType == CheckTypeAll {
						log.With("hash", video.SourceHash, "relate", video.Bangumi, "type", "video").Info("unpin")
					}
				}

			}
			if !seed.SkipTypeVerify(model.TypeSlice, p.skip...) {
				if _, b := pinned[video.M3U8Hash]; b {
					if p.checkType == CheckTypePin || p.checkType == CheckTypeAll {
						log.With("hash", video.M3U8Hash, "relate", video.Bangumi, "type", "slice").Info("pinned")

					}
					pinned[video.M3U8Hash] = video
				} else {
					if p.checkType == CheckTypeUnpin || p.checkType == CheckTypeAll {
						log.With("hash", video.M3U8Hash, "relate", video.Bangumi, "type", "slice").Info("unpin")
					}
				}

			}
		}
	}
	log.Info("pin done")

	for hash, v := range pinned {
		if v != nil {
			log.With("hash", hash, "relate", v.Bangumi).Info("add pin")
			p := &model.Pin{
				PinHash: hash,
				PeerID:  myid.ID,
				//VideoID: "",
			}

			e = a.PushTo(seed.DatabaseCallback(p, func(database *seed.Database, eng *xorm.Engine, v interface{}) (e error) {
				p := v.(*model.Pin)
				return model.AddOrUpdatePin(eng.NoCache(), p)
			}))
			if e != nil {
				log.Error(e)
				break
			}
		}
	}

	close(v)
}

func (p *pinCheck) pinUnfinishedCall(a *seed.API, api *httpapi.HttpApi) {
	u := make(chan *model.Unfinished)
	e := a.PushTo(seed.DatabaseUnfinishedCall(u, func(session *xorm.Session) *xorm.Session {
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
	log.With("total", len(pinned)).Info("pinned")

	myid, e := seed.MyID(a)
	if e != nil {
		return
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
	log.Info("pin done")

	for _, u := range pinned {
		if u != nil {
			log.With("hash", u.Hash, "relate", u.Relate, "type", u.Type).Info("add pin")
			p := &model.Pin{
				PinHash: u.Hash,
				PeerID:  myid.ID,
				//VideoID: "",
			}

			e = a.PushTo(seed.DatabaseCallback(p, func(database *seed.Database, eng *xorm.Engine, v interface{}) (e error) {
				p := v.(*model.Pin)
				return model.AddOrUpdatePin(eng.NoCache(), p)
			}))
			if e != nil {
				log.Error(e)
				break
			}
		}
	}

	close(u)
}

type pinSync struct {
	from  string
	skip  []interface{}
	table PinTable
}

func (p *pinSync) Call(a *seed.API, api *httpapi.HttpApi) error {
	if p.from == "" {
		return nil
	}
	ma, err := multiaddr.NewMultiaddr(p.from)
	if err != nil {
		return err

	}
	pi, err := peer.AddrInfoFromP2pAddr(ma)
	if err != nil {
		return err
	}
	err = api.Swarm().Connect(a.Context(), *pi)
	if err != nil {
		return err
	}
	log.Info("pin sync")
	switch p.table {
	case PinTableUnfinished:
		p.pinUnfinishedCall(a, api)
	case PinTableVideo:
		p.pinVideoCall(a, api)
	case PinTablePin:
		p.pinPinCall(a, api)
	}

	return nil
}

func (p *pinSync) pinVideoCall(a *seed.API, api *httpapi.HttpApi) {
	pp := make(chan *model.Pin)
	err := a.PushTo(seed.DatabasePinCall(pp, func(session *xorm.Session) *xorm.Session {
		idx := strings.LastIndex(p.from, "/") + 1
		from := p.from
		if idx >= 0 {
			from = p.from[idx:]
		}
		return session.Where("peer_id = ?", from)
	}))
	if err != nil {
		return
	}

	var pins []*model.Pin
ChanEnd:
	for {
		select {
		case pin := <-pp:
			if pin == nil {
				break ChanEnd
			}
			pins = append(pins, pin)
		}
	}

	for _, pin := range pins {
		err := a.PushTo(seed.DatabaseCallback(pin, func(database *seed.Database, eng *xorm.Engine, v interface{}) (e error) {
			pin := v.(*model.Pin)
			session := eng.NoCache()
			if !seed.SkipTypeVerify(model.TypeSlice, p.skip...) {
				session = session.Or("m3u8_hash = ?", pin.PinHash)
			}
			if !seed.SkipTypeVerify(model.TypeVideo, p.skip...) {
				session = session.Or("source_hash = ?", pin.PinHash)
			}
			if !seed.SkipTypeVerify(model.TypePoster, p.skip...) {
				session = session.Or("poster_hash = ?", pin.PinHash)
			}
			if !seed.SkipTypeVerify(model.TypeThumb, p.skip...) {
				session = session.Or("thumb_hash = ?", pin.PinHash)
			}
			vs := new([]*model.Video)
			i, e := session.FindAndCount(vs)
			if e != nil {
				return e
			}
			if i > 0 {
				log.With("hash", pin.PinHash, "peer_id", pin.PeerID, "video", (*vs)[0].Bangumi).Info("pinning")
				err := api.Pin().Add(a.Context(), path.New(pin.PinHash), func(settings *options.PinAddSettings) error {
					settings.Recursive = true
					return nil
				})
				if err != nil {
					log.Error(err)
				}
			}
			return nil
		}))
		if err != nil {
			log.Error(err)
		}
	}

}

func (p *pinSync) pinUnfinishedCall(a *seed.API, api *httpapi.HttpApi) {
	pp := make(chan *model.Pin)
	err := a.PushTo(seed.DatabasePinCall(pp, func(session *xorm.Session) *xorm.Session {
		idx := strings.LastIndex(p.from, "/") + 1
		from := p.from
		if idx >= 0 {
			from = p.from[idx:]
		}
		return session.Where("peer_id = ?", from)
	}))
	if err != nil {
		return
	}
	var pins []*model.Pin
ChanEnd:
	for {
		select {
		case pin := <-pp:
			if pin == nil {
				break ChanEnd
			}
			pins = append(pins, pin)
		}
	}

	for _, pin := range pins {
		err := a.PushTo(seed.DatabaseCallback(pin, func(database *seed.Database, eng *xorm.Engine, v interface{}) (e error) {
			pin := v.(*model.Pin)
			session := eng.NoCache()
			if !seed.SkipTypeVerify(model.TypeSlice, p.skip...) {
				session = session.Or("type = ?", model.TypeSlice)
			}
			if !seed.SkipTypeVerify(model.TypeVideo, p.skip...) {
				session = session.Or("type = ?", model.TypeVideo)
			}
			if !seed.SkipTypeVerify(model.TypePoster, p.skip...) {
				session = session.Or("type = ?", model.TypePoster)
			}
			if !seed.SkipTypeVerify(model.TypeThumb, p.skip...) {
				session = session.Or("type = ?", model.TypeThumb)
			}
			us := new([]*model.Unfinished)
			i, e := session.And("hash = ?", pin.PinHash).FindAndCount(us)
			if e != nil {
				return e
			}
			if i > 0 {
				log.With("hash", pin.PinHash, "peer_id", pin.PeerID, "type", (*us)[0].Type, "relate", (*us)[0].Relate).Info("pinning")
				err := api.Pin().Add(a.Context(), path.New(pin.PinHash), func(settings *options.PinAddSettings) error {
					settings.Recursive = true
					return nil
				})
				if err != nil {
					log.Error(err)
				}
			}
			return nil

		}))
		if err != nil {
			log.Error(err)
		}
	}

}

func (p *pinSync) pinPinCall(a *seed.API, api *httpapi.HttpApi) {
	pp := make(chan *model.Pin)
	err := a.PushTo(seed.DatabasePinCall(pp, func(session *xorm.Session) *xorm.Session {
		idx := strings.LastIndex(p.from, "/") + 1
		from := p.from
		if idx >= 0 {
			from = p.from[idx:]
		}
		log.With("from", from).Info("sync")
		return session.Where("peer_id = ?", from)
	}))
	if err != nil {
		log.Error(err)
	}
ChanEnd:
	for {
		select {
		case pin := <-pp:
			if pin == nil {
				break ChanEnd
			}
			log.With("hash", pin.PinHash, "peer_id", pin.PeerID).Info("pinning")
			err := api.Pin().Add(a.Context(), path.New(pin.PinHash), func(settings *options.PinAddSettings) error {
				settings.Recursive = true
				return nil
			})
			if err != nil {
				log.Error(err)
			}
		}
	}

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
