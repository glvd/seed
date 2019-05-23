package seed

import (
	"github.com/yinhevr/seed/model"
	"golang.org/x/xerrors"
	"sync"
)

func pin(wg *sync.WaitGroup, hash string) {
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
	SwarmAddList(video.SourcePeerList)
	log.Info("pin video:", video.Bangumi)
	wg.Add(1)
	//log.Info("pin poster:", video.Poster)
	go pin(wg, video.Poster)
	if poster {
		return
	}
	for _, value := range video.VideoGroupList {
		log.Infof("list:%+v", value)
		for _, val := range value.Object {
			//log.Info("pin media:", val.Link.Hash)
			wg.Add(1)
			go pin(wg, val.Link.Hash)
		}
	}
}

// QuickPin ...
func QuickPin(checksum string, check bool) (e error) {
	log.Info("pin checksum:", checksum)
	var uncategorizeds []*model.Uncategorized
	if checksum == "" {
		uncategorizeds, e = model.AllUncategorized(check)
		if e != nil {
			return e
		}

	} else {
		uncategorized, e := model.FindUncategorized(checksum, check)
		if e != nil {
			return e
		}
		uncategorizeds = append(uncategorizeds, uncategorized)
	}
	if !check {
		return
	}
	for _, v := range uncategorizeds {
		pin(nil, v.Hash)
		v.Sync = true
		i, e := model.DB().Cols("sync").Update(v)
		if e != nil {
			return xerrors.Errorf("uncategorized nothing updated with:%d,%+v", i, e)
		}
	}

	return nil
}

// Pin ...
func Pin(ban string, poster, check bool) (e error) {
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
