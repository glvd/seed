package model

import (
	"github.com/go-xorm/xorm"
	"golang.org/x/xerrors"
)

// Type ...
type Type string

// TypeOther ...
const TypeOther Type = "other"

// TypeVideo ...
const TypeVideo Type = "video"

// TypeSlice ...
const TypeSlice Type = "slice"

// TypePoster ...
const TypePoster Type = "poster"

// TypeThumb ...
const TypeThumb Type = "thumb"

// Unfinished 未分类
type Unfinished struct {
	Model       `xorm:"extends"`
	Checksum    string       `xorm:"default() checksum"`            //sum值
	Type        Type         `xorm:"default() type"`                //类型
	Relate      string       `xorm:"default()" json:"relate"`       //关联信息
	Name        string       `xorm:"default() name"`                //名称
	Hash        string       `xorm:"default() hash"`                //哈希地址
	IsVideo     bool         `xorm:"default(0)"`                    //视频文件
	Sharpness   string       `xorm:"default()" json:"sharpness"`    //清晰度
	Caption     string       `xorm:"default()" json:"caption"`      //字幕
	Encrypt     bool         `json:"encrypt"`                       //加密
	Key         string       `xorm:"default()"json:"key"`           //秘钥
	M3U8        string       `xorm:"m3u8 default()" json:"m3u8"`    //M3U8名
	SegmentFile string       `xorm:"default()" json:"segment_file"` //ts切片名
	Sync        bool         `xorm:"notnull default(0)"`            //是否已同步
	Object      *VideoObject `xorm:"json" json:"object,omitempty"`  //视频信息
}

func init() {
	RegisterTable(Unfinished{})
}

// AllUnfinished ...
func AllUnfinished(session *xorm.Session, limit int, start ...int) (unfins *[]*Unfinished, e error) {
	unfins = new([]*Unfinished)
	session = MustSession(session)
	if limit > 0 {
		session = session.Limit(limit, start...)
	}
	if err := session.Find(unfins); err != nil {
		return nil, err
	}
	return unfins, nil
}

// FindUnfinished ...
func FindUnfinished(session *xorm.Session, checksum string) (unfin *Unfinished, e error) {
	unfin = new(Unfinished)
	b, e := MustSession(session).Where("checksum = ?", checksum).Get(unfin)
	if e != nil || !b {
		return nil, xerrors.New("uncategorize not found!")
	}
	return unfin, nil
}

// AddOrUpdateUnfinished ...
func AddOrUpdateUnfinished(unfin *Unfinished) (e error) {
	tmp := new(Unfinished)
	found, e := DB().Table(unfin).Where("checksum = ?", unfin.Checksum).Get(tmp)
	if e != nil {
		return e
	}
	if found {
		unfin.Version = tmp.Version
		unfin.ID = tmp.ID
		_, e = DB().ID(unfin.ID).Update(unfin)
		return
	}
	_, e = DB().InsertOne(unfin)
	return
}
