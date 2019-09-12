package model

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-xorm/xorm"
	"github.com/ipfs/interface-go-ipfs-core/path"
)

// Pin ...
type Pin struct {
	Model    `xorm:"extends" json:"-"`
	PinHash  string `xorm:"pin_hash"`
	PeerID   string `xorm:"peer_id"`
	VideoID  string `xorm:"video_id"`
	IsPinned string `xorm:"is_pinned"` //TODO:status is now not used
}

func init() {
	RegisterTable(Pin{})
}

//AllPin find pins
func AllPin(session *xorm.Session, limit int, start ...int) (pins *[]*Pin, e error) {
	pins = new([]*Pin)
	session = MustSession(session)
	if limit > 0 {
		session = session.Limit(limit, start...)
	}
	if err := session.Find(pins); err != nil {
		return nil, err
	}
	return pins, nil
}

//FindPin find one pin
func FindPin(session *xorm.Session, ph string) (pin *Pin, e error) {
	pin = new(Pin)
	b, e := MustSession(session).Where("pin_hash = ?", ph).Get(pin)
	if e != nil || !b {
		return nil, errors.New("pin not found")
	}
	return pin, nil
}

// PinHash ...
func PinHash(path path.Resolved) string {
	ss := strings.Split(path.String(), "/")
	if len(ss) == 3 {
		return ss[2]
	}
	return ""
}

// UpdatePinVideoID ...
func UpdatePinVideoID(session *xorm.Session, p *Pin) (e error) {
	videos := new([]Video)
	i, e := session.Clone().Table(&Video{}).Where("m3u8_hash = ?", p.PinHash).
		Or("source_hash = ?", p.PinHash).
		Or("thumb_hash = ?", p.PinHash).
		Or("poster_hash = ?", p.PinHash).FindAndCount(videos)
	if e != nil {
		return e
	}

	if i == 1 {
		p.VideoID = (*videos)[0].ID
	} else if i > 1 {
		p.VideoID = fmt.Sprintf("ids(%d)", i)
	} else {
		p.VideoID = "dummy"
	}
	return AddOrUpdatePin(session.Clone(), p)
}

// AddOrUpdatePin ...
func AddOrUpdatePin(session *xorm.Session, p *Pin) (e error) {
	tmp := new(Pin)
	var found bool
	if p.ID != "" {
		found, e = session.Clone().ID(p.ID).Get(tmp)
	} else {
		found, e = session.Clone().Where("pin_hash = ?", p.PinHash).
			And("peer_id = ?", p.PeerID).Get(tmp)
	}
	if e != nil {
		return e
	}
	if found {
		//only slice need update,video update for check
		log.With("peer_id", p.PeerID, "pin_hash", p.PinHash).Info("exist")
		return
	}
	_, e = session.Clone().InsertOne(p)
	return
}

// IsExist ...
//func (p *Pin) IsExist() bool {
//	i, e := DB().Table(&Pin{}).Where("pin_hash = ?", p.PinHash).Count()
//	log.With("pin_hash", p.PinHash, "num", i).Info("check exist")
//	if e != nil || i <= 0 {
//		return false
//	}
//	return true
//}
