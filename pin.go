package seed

import (
	"context"
	shell "github.com/godcong/go-ipfs-restapi"
	"sync"

	"github.com/yinhevr/seed/model"
	"golang.org/x/xerrors"
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
				wg.Add(1)
				go pinHash(wg, v.Hash)
				wg.Add(1)
				go pinHash(wg, v.SliceHash)
				wg.Wait()
			default:
				//nothing to do
			}
		}
	}
}

func pinHash(wg *sync.WaitGroup, hash string) {
	log.Info("pin:", hash)
	e := rest.Pin(hash)
	if e != nil {
		log.Error("pin error:", hash, e)
		return
	}
	if wg != nil {
		wg.Done()
	}

	log.Info("pinned:", hash)
}

func pinVideo(wg *sync.WaitGroup, poster bool, video *model.Video) {
	//SwarmAddList(video.SourcePeerList)
	log.Info("pin video:", video.Bangumi)
	wg.Add(1)
	//log.Info("pin poster:", video.Poster)
	go pinHash(wg, video.PosterHash)
	if poster {
		return
	}
	//for _, value := range video.VideoGroupList {
	//	log.Infof("list:%+v", value)
	//	for _, val := range value.Object {
	//log.Info("pin media:", val.Link.Hash)
	//wg.Add(1)
	//go pin(wg, val.Link.Hash)
	//}
	//}
}

// QuickPin ...
func QuickPin(checksum string, check bool) (e error) {
	log.Info("pin checksum:", checksum)
	session := model.DB().Where("sync = ?", !check)

	var unfins []*model.Unfinished
	if checksum != "" {
		session = session.Where("checksum = ?", checksum)

	}
	unfins, e = model.AllUnfinished(session, 0, 500)
	if e != nil {
		return e
	}

	for _, v := range unfins {
		pinHash(nil, v.Hash)
	}
	//, func(hash string) {
	//	v.Sync = true
	//	i, e := model.DB().Cols("sync").Update(v)
	//	if e != nil {
	//		log.Errorf("Unfinished nothing updated with:%d,%+v", i, e)
	//	}
	//}
	return nil
}

// PinProc ...
func PinProc(ban string, poster, check bool) (e error) {
	wg := sync.WaitGroup{}
	var videos []*model.Video
	if ban == "" {
		videos, e = model.AllVideos(check)
		if e != nil {
			return e
		}

	} else {
		videos = append(videos, new(model.Video))
		b, e := model.FindVideo(ban, videos[0], check)
		if e != nil || !b {
			return xerrors.Errorf("nothing updated with:%t,%+v", b, e)
		}
	}
	for _, video := range videos {
		pinVideo(&wg, poster, video)
	}
	wg.Wait()
	if !check {
		return
	}
	for _, video := range videos {
		video.Sync = true
		i, e := model.DB().Cols("sync").Update(video)
		if e != nil {
			return xerrors.Errorf("video nothing updated with:%d,%+v", i, e)
		}
	}

	log.Info("success")
	return nil
}
