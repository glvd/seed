package model

import (
	"fmt"
	"strings"

	"github.com/go-xorm/xorm"
	"github.com/ipfs/interface-go-ipfs-core/path"
	"golang.org/x/xerrors"
)

// Pin ...
type Pin struct {
	Model   `xorm:"extends" json:"-"`
	PinHash string   `xorm:"pin_hash"`
	PeerID  []string `xorm:"machine_id"`
	VideoID string   `xorm:"video_id"`
}

func init() {
	RegisterTable(Pin{})
}

// PinHash ...
func PinHash(path path.Resolved) string {
	ss := strings.Split(path.String(), "/")
	if len(ss) == 3 {
		return ss[2]
	}
	return ""
}

// UpdateVideo ...
func (p *Pin) UpdateVideo() (e error) {
	return updatePinVideoID(DB().NewSession(), p)
}

func updatePinVideoID(session *xorm.Session, p *Pin) (e error) {
	videos := new([]Video)
	i, e := session.Table(&Video{}).Where("m3u8_hash = ?", p.PinHash).
		Or("source_hash = ?", p.PinHash).
		Or("thumb_hash = ?", p.PinHash).
		Or("poster_hash = ?", p.PinHash).FindAndCount(videos)
	if e != nil {
		return e
	}
	if i <= 0 {
		return xerrors.New("nothing found")
	}
	if i == 1 {
		p.VideoID = (*videos)[0].ID
	} else if i > 1 {
		p.VideoID = fmt.Sprintf("ids(%d)", i)
	}
	return AddOrUpdatePin(p)
}

// AddOrUpdatePin ...
func AddOrUpdatePin(p *Pin) (e error) {
	tmp := new(Pin)
	var found bool
	if p.ID != "" {
		found, e = DB().ID(p.ID).Get(tmp)
	} else {
		found, e = DB().Where("pin_hash = ?", p.PinHash).Get(tmp)
	}
	if e != nil {
		return e
	}
	if found {
		//only slice need update,video update for check
		p.Version = tmp.Version
		p.ID = tmp.ID
		ids := tmp.PeerID
		for _, pid := range p.PeerID {
			for _, tid := range tmp.PeerID {
				if pid == tid {
					pid = ""
					break
				}
			}
			if pid != "" {
				ids = append(ids, pid)
			}
		}
		p.PeerID = ids
		_, e = DB().ID(p.ID).Update(p)
		return
	}
	_, e = DB().InsertOne(p)
	return
}

// IsExist ...
func (p *Pin) IsExist() bool {
	i, e := DB().Table(&Pin{}).Where("pin_hash = ?", p.PinHash).Count()
	log.With("pin_hash", p.PinHash, "num", i).Info("check exist")
	if e != nil || i <= 0 {
		return false
	}
	return true
}
