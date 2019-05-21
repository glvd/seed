package model

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"golang.org/x/xerrors"
	"io"
	"os"
)

// Uncategorized 未分类
type Uncategorized struct {
	Model    `xorm:"extends"`
	Checksum string
	Type     string
	Name     string
	Hash     string
	IsVideo  bool           `xorm:"notnull default(0)"`
	Sync     bool           `xorm:"notnull default(0)"`
	Object   []*VideoObject `xorm:"json" json:"object,omitempty"` //视频信息
}

func init() {
	RegisterTable(Uncategorized{})
}

// AllUncategorized ...
func AllUncategorized(check bool) ([]*Uncategorized, error) {
	var uncats []*Uncategorized
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

// FindUncategorized ...
func FindUncategorized(checksum string, check bool) (*Uncategorized, error) {
	var uncat Uncategorized
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

// AddOrUpdateUncategorized ...
func AddOrUpdateUncategorized(uncat *Uncategorized) (e error) {
	log.Infof("%+v", *uncat)
	i, e := DB().Table(uncat).Where("checksum = ?", uncat.Checksum).And("type = ?", uncat.Type).Count()
	if e != nil {
		return e
	}
	if i > 0 {
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
