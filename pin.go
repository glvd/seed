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

func pinVideo(wg *sync.WaitGroup, poster bool, video *model.Video) {
	go swarmConnects(video.SourcePeerList)
	logrus.Info("pin video:", video.Bangumi)
	wg.Add(1)
	//logrus.Info("pin poster:", video.Poster)
	go pin(wg, video.Poster)
	if poster {
		return
	}
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
func QuickPin(checksum string, check bool) (e error) {
	logrus.Info("pin checksum:", checksum)
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
	for _, v := range uncategorizeds {
		logrus.Info("pin:", v.Hash)
		pin(nil, v.Hash)
		v.Sync = true
		i, e := model.DB().Update(v)
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

	for _, video := range videos {
		video.Sync = true
		i, e := model.DB().Update(video)
		if e != nil {
			return xerrors.Errorf("video nothing updated with:%d,%+v", i, e)
		}
	}

	logrus.Info("success")
	return nil
}
