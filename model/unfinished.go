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
	SliceHash   string       `xorm:"default()" json:"slice_hash"`   //切片HASH
	IsVideo     bool         `xorm:"default(0)"`                    //视频文件
	Sharpness   string       `xorm:"default()" json:"sharpness"`    //清晰度
	Sliced      bool         `json:"sliced"`                        //切片
	Encrypt     bool         `json:"encrypt"`                       //加密
	Key         string       `xorm:"default()"json:"key"`           //秘钥
	M3U8        string       `xorm:"m3u8 default()" json:"m3u8"`    //M3U8名
	Caption     string       `xorm:"default()" json:"caption"`      //字幕
	Sync        bool         `xorm:"notnull default(0)"`            //是否已同步
	SegmentFile string       `xorm:"default()" json:"segment_file"` //ts切片名
	Object      *VideoObject `xorm:"json" json:"object,omitempty"`  //视频信息
	SliceObject *VideoObject `xorm:"json" json:"slice_object,omitempty"`
}

func init() {
	RegisterTable(Unfinished{})
}

// AllUnfinished ...
func AllUnfinished(session *xorm.Session, limit int, start ...int) (unfins []*Unfinished, e error) {
	unfins = []*Unfinished{}
	if err := MustSession(session).Limit(limit, start...).Find(&unfins); err != nil {
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
		if _, e = DB().Where("checksum = ?", unfin.Checksum).Update(unfin); e != nil {
			return e
		}
		return nil
	}
	_, e = DB().InsertOne(unfin)
	return
}
