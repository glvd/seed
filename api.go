package seed

import (
	"context"
	"github.com/godcong/go-ipfs-restapi"
	"github.com/yinhevr/seed/model"
	"golang.org/x/xerrors"
	"strings"
	"sync"
	"time"
)

// QuickConnect ...
func QuickConnect(shell *shell.Shell, addr string) {
	var e error
	go func() {
		for {
			e = SwarmConnect(shell, addr)
			if e != nil {
				return
			}
			time.Sleep(30 * time.Second)
		}
	}()
}

var swarms = sync.Pool{}

// PoolSwarmConnect ...
func PoolSwarmConnect(shell *shell.Shell) {
	log.Info("PoolSwarmConnect running")
	for {
		if s := swarms.Get(); s != nil {
			sp, b := s.(*model.SourcePeer)
			if b {
				e := SwarmConnect(shell, swarmAddress(sp))
				log.Info(swarmAddress(sp))
				if e != nil {
					log.Error("swarm connect err:", e)
				}
				swarms.Put(sp)
			}
			time.Sleep(30 * time.Second)
			continue
		}
		time.Sleep(5 * time.Second)
	}
}

// SwarmAdd ...
func SwarmAdd(sp *model.SourcePeer) {
	swarms.Put(sp)
}

// SwarmAddAddress ...
func SwarmAddAddress(addr string) {
	swarms.Put(AddressSwarm(addr))
}

// SwarmAddList ...
func SwarmAddList(sps []*model.SourcePeer) {
	for _, v := range sps {
		SwarmAdd(v)
	}
}

// SwarmConnect ...
func SwarmConnect(shell *shell.Shell, addr string) (e error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	log.Info("connect to:", addr)
	if err := shell.SwarmConnect(ctx, addr); err != nil {
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

// AddressSwarm ...
func AddressSwarm(address string) (peer *model.SourcePeer) {
	ss := strings.Split(address, "/")
	size := len(ss)
	log.Info("address:", address)
	log.Info("size:", size)
	if size < 7 {
		return &model.SourcePeer{}
	}
	return &model.SourcePeer{
		SourcePeerDetail: model.SourcePeerDetail{
			Addr: strings.Join(ss[:size-2], "/"),
			Peer: ss[size-1],
		},
	}
}

func swarmConnectTo(shell *shell.Shell, peer *model.SourcePeer) (e error) {
	address := swarmAddress(peer)
	if address == "" {
		return xerrors.New("null address")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	log.Info("connect to:", address)
	if err := shell.SwarmConnect(ctx, address); err != nil {
		return err
	}
	return
}

func swarmConnects(shell *shell.Shell, peers []*model.SourcePeer) {
	if peers == nil {
		return
	}

	var nextPeers []*model.SourcePeer
	for _, value := range peers {
		e := swarmConnectTo(shell, value)
		if e != nil {
			//log.Error(e)
			time.Sleep(30 * time.Second)
			continue
		}
		//filter the error peers
		nextPeers = append(nextPeers, value)
		time.Sleep(30 * time.Second)
	}
	//rerun when connect is end
	swarmConnects(shell, nextPeers)
}
