package model

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"golang.org/x/xerrors"
	"io"
	"os"
)

// Unfinished 未分类
type Unfinished struct {
	Model       `xorm:"extends"`
	Checksum    string       `xorm:"checksum default()"`
	Type        string       `xorm:"type default()"`
	Name        string       `xorm:"default()" xorm:"name"`
	Hash        string       `xorm:"default()" xorm:"hash"`         //哈希地址
	IsVideo     bool         `xorm:"default(0)"`                    //视频文件
	Sharpness   string       `xorm:"default()" json:"sharpness"`    //清晰度
	Sync        bool         `xorm:"default(0)"`                    //是否同步
	Sliced      bool         `json:"sliced"`                        //切片
	Encrypt     bool         `json:"encrypt"`                       //加密
	Key         string       `xorm:"default()"json:"key"`           //秘钥
	M3U8        string       `xorm:"m3u8 default()" json:"m3u8"`    //M3U8名
	Caption     string       `xorm:"default()" json:"caption"`      //字幕
	SegmentFile string       `xorm:"default()" json:"segment_file"` //ts切片名
	Object      *VideoObject `xorm:"json" json:"object,omitempty"`  //视频信息
}

func init() {
	RegisterTable(Unfinished{})
}

// AllUnfinished ...
func AllUnfinished(check bool) ([]*Unfinished, error) {
	var uncats []*Unfinished
	if check {
		if err := DB().Where("sync = ?", !check).Find(&uncats); err != nil {
			return nil, err
		}
	} else {
		if err := DB().Find(&uncats); err != nil {
			return nil, err
		}
	}
	return uncats, nil
}

// FindUnfinished ...
func FindUnfinished(checksum string, check bool) (*Unfinished, error) {
	var uncat Unfinished
	if check {
		b, e := DB().Where("type = ?", "m3u8").Where("sync = ?", !check).Where("checksum = ?", checksum).Get(&uncat)
		if e != nil || !b {
			return nil, xerrors.New("uncategorize not found!")
		}
	} else {
		b, e := DB().Where("type = ?", "m3u8").Where("checksum = ?", checksum).Get(&uncat)
		if e != nil || !b {
			return nil, xerrors.New("uncategorize not found!")
		}
	}
	return &uncat, nil
}

// AddOrUpdateUnfinished ...
func AddOrUpdateUnfinished(uncat *Unfinished) (e error) {
	log.Infof("%+v", *uncat)
	tmp := new(Unfinished)
	b, e := DB().Table(uncat).Where("checksum = ?", uncat.Checksum).And("type = ?", uncat.Type).Get(tmp)
	if e != nil {
		return e
	}
	if b {
		uncat.Version = tmp.Version
		if _, err := DB().Where("checksum = ?", uncat.Checksum).Update(uncat); err != nil {
			return err
		}
		return nil
	}
	if _, err := DB().InsertOne(uncat); err != nil {
		return err
	}
	return nil
}

// Checksum ...
func Checksum(filepath string) string {
	hash := md5.New()
	file, e := os.OpenFile(filepath, os.O_RDONLY, os.ModePerm)
	if e != nil {
		return ""
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	_, e = io.Copy(hash, reader)
	if e != nil {
		return ""
	}
	return fmt.Sprintf("%x", hash.Sum(nil))
}
