package seed

import (
	"github.com/girlvr/seed/model"
	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
	"sync"
)

// Pin ...
func Pin(ban string) (e error) {
	wg := sync.WaitGroup{}
	if ban == "" {
		v, e := model.AllVideos()
		if e != nil {
			return e
		}
		for _, videos := range v {
			for _, value := range videos.VideoGroupList {
				for _, val := range value.Object {
					wg.Add(1)
					go func(hash string) {
						e := rest.Pin(hash)
						if e != nil {
							logrus.Error(e)
						}
						wg.Done()
					}(val.Link.Hash)
				}
			}
		}
	} else {
		var video model.Video
		b, e := model.FindVideo(ban, &video)
		if e != nil || !b {
			return xerrors.New("nil video")
		}

		for _, value := range video.VideoGroupList {
			for _, val := range value.Object {
				wg.Add(1)
				go func(hash string) {
					e := rest.Pin(hash)
					if e != nil {
						logrus.Error(e)
					}
					wg.Done()
				}(val.Link.Hash)
			}
		}
	}
	wg.Wait()
	logrus.Info("success")
	return nil
}
