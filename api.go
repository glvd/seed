package seed

import (
	"context"
	"github.com/godcong/go-ipfs-restapi"
	"github.com/sirupsen/logrus"
	"github.com/yinhevr/seed/model"
	"golang.org/x/xerrors"
	"sync"
	"time"
)

var rest *shell.Shell

// InitShell ...
func InitShell(s string) {
	logrus.Info("ipfs shell:", s)
	rest = shell.NewShell(s)
}

// QuickConnect ...
func QuickConnect(addr string) {
	var e error
	go func() {
		for {
			e = SwarmConnect(addr)
			if e != nil {
				return
			}
			time.Sleep(30 * time.Second)
		}
	}()
}

var swarms = sync.Pool{}

// PoolSwarmConnect ...
func PoolSwarmConnect() {
	SwarmAdd(&model.SourcePeer{
		SourcePeerDetail: &model.SourcePeerDetail{
			Addr: "/ip4/47.101.169.94/tcp/4001",
			Peer: "Qmcoz66NZhcegp58st53Khsd2mgqnkLojQx7mtjAA3EPCS",
		},
	})
	for {
		if s := swarms.Get(); s != nil {
			sp, b := s.(*model.SourcePeer)
			if b {
				e := SwarmConnect(swarmAddress(sp))
				if e != nil {
					time.Sleep(30 * time.Second)
					continue
				}
				swarms.Put(sp)
			}
		}
		time.Sleep(30 * time.Second)
	}
}

// SwarmAdd ...
func SwarmAdd(sp *model.SourcePeer) {
	swarms.Put(sp)
}

// SwarmAddList ...
func SwarmAddList(sps []*model.SourcePeer) {
	for _, v := range sps {
		SwarmAdd(v)
	}
}

// SwarmConnect ...
func SwarmConnect(addr string) (e error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	logrus.Info("connect to:", addr)
	if err := rest.SwarmConnect(ctx, addr); err != nil {
		cancel()
		return err
	}
	return
}
func swarmAddress(peer *model.SourcePeer) string {
	if peer != nil {
		return peer.Addr + "/ipfs/" + peer.Peer
	}
	return ""
}

func swarmConnectTo(peer *model.SourcePeer) (e error) {
	address := swarmAddress(peer)
	if address == "" {
		return xerrors.New("null address")
	}
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
			//logrus.Error(e)
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
