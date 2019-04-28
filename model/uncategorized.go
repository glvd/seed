package model

import (
	"bufio"
	"crypto/md5"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
)

// Uncategorized 未分类
type Uncategorized struct {
	Model    `xorm:"extends"`
	Checksum string
	Name     string
	Hash     string
	IsVideo  bool
	Object   []*VideoObject `xorm:"json" json:"object,omitempty"` //视频信息
}

func init() {
	RegisterTable(Uncategorized{})
}

// AddOrUpdateVideo ...
func AddOrUpdateUncategorized(uncat *Uncategorized) (e error) {
	log.Printf("%+v", *uncat)
	if uncat.ID != "" {
		if _, err := DB().Where("name = ?", uncat.Name).Update(uncat); err != nil {
			return err
		}
		return nil
	}
	if _, err := DB().InsertOne(uncat); err != nil {
		return err
	}
	return nil
}

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
