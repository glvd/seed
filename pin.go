package seed

import (
	"context"
	"github.com/girlvr/seed/model"
	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
	"sync"
	"time"
)

func swarmConnectTo(peer *model.SourcePeer) (e error) {
	address := peer.Addr + "/" + peer.Peer
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
	if err := rest.SwarmConnect(ctx, address); err != nil {
		return err
	}
	return
}

func swarmConnects(peers []*model.SourcePeer) {
	if peers == nil {
		return
	}

	for idx, value := range peers {
		e := swarmConnectTo(value)
		if e != nil {
			logrus.Error(e)
			//filter the error peers
			peers = append(peers[:idx], peers[idx+1:]...)
		}
		time.Sleep(30 * time.Second)
	}
	//rerun when connect is end
	swarmConnects(peers)
}

func pin(wg *sync.WaitGroup, hash string) {
	e := rest.Pin(hash)
	if e != nil {
		logrus.Error(e)
	}
	wg.Done()
}

func pinVideo(wg *sync.WaitGroup, video *model.Video) {
	go swarmConnects(video.SourcePeerList)
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
