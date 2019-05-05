package seed

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/yinhevr/seed/model"
	"golang.org/x/xerrors"
	"sync"
	"time"
)

func swarmConnectTo(peer *model.SourcePeer) (e error) {
	address := peer.Addr + "/ipfs/" + peer.Peer
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	logrus.Info("connect to:", address)
	if err := rest.SwarmConnect(ctx, address); err != nil {
		cancel()
		return err
	}
	return
}

func swarmConnects(peers []*model.SourcePeer) {
	if peers == nil {
		return
	}

	var nextPeers []*model.SourcePeer

	for _, value := range peers {
		e := swarmConnectTo(value)
		if e != nil {
			logrus.Error(e)
			time.Sleep(30 * time.Second)
			continue
		}
		//filter the error peers
		nextPeers = append(nextPeers, value)
		time.Sleep(30 * time.Second)
	}
	//rerun when connect is end
	swarmConnects(nextPeers)
}

func pin(wg *sync.WaitGroup, hash string) {
	logrus.Info("pin:", hash)
	e := rest.Pin(hash)
	if e != nil {
		logrus.Error(e)
	}
	if wg != nil {
		wg.Done()
	}
}

func pinVideo(wg *sync.WaitGroup, video *model.Video) {
	go swarmConnects(video.SourcePeerList)
	logrus.Info("pin video:", video.Bangumi)
	wg.Add(1)
	//logrus.Info("pin poster:", video.Poster)
	go pin(wg, video.Poster)
	for _, value := range video.VideoGroupList {
		logrus.Infof("list:%+v", value)
		for _, val := range value.Object {
			//logrus.Info("pin media:", val.Link.Hash)
			wg.Add(1)
			go pin(wg, val.Link.Hash)
		}
	}
}

// QuickPin ...
func QuickPin(checksum string) (e error) {
	logrus.Info("pin checksum:", checksum)
	if checksum == "" {
		uncategorizeds, e := model.AllUncategorized()
		if e != nil {
			return e
		}
		for _, v := range uncategorizeds {

			pin(nil, v.Hash)
		}
	}
	uncategorized, e := model.FindUncategorized(checksum)
	if e != nil {
		return e
	}
	pin(nil, uncategorized.Hash)
	return nil
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
