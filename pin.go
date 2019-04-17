package seed

import (
	"github.com/girlvr/seed/model"
	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
	"sync"
)

func connectSwarm() {

}

func pin(wg *sync.WaitGroup, hash string) {
	e := rest.Pin(hash)
	if e != nil {
		logrus.Error(e)
	}
	wg.Done()
}

func pinVideo(wg *sync.WaitGroup, video *model.Video) {
	logrus.Info("pin video:", video.Bangumi)
	wg.Add(1)
	logrus.Info("pin poster:", video.Poster)
	go pin(wg, video.Poster)
	for _, value := range video.VideoGroupList {
		logrus.Infof("list:%+v", value)
		for _, val := range value.Object {
			logrus.Info("pin media:", val.Link.Hash)
			wg.Add(1)
			go pin(wg, val.Link.Hash)
		}
	}
}

// Pin ...
func Pin(ban string) (e error) {
	wg := sync.WaitGroup{}
	if ban == "" {
		v, e := model.AllVideos()
		if e != nil {
			return e
		}
		for _, video := range v {
			pinVideo(&wg, video)
		}
	} else {
		var video model.Video
		b, e := model.FindVideo(ban, &video)
		if e != nil || !b {
			return xerrors.New("nil video")
		}
		pinVideo(&wg, &video)
	}
	wg.Wait()
	logrus.Info("success")
	return nil
}
