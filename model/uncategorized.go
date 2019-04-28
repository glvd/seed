package model

import (
	log "github.com/sirupsen/logrus"
)

// Uncategorized 未分类
type Uncategorized struct {
	Model   `xorm:"extends"`
	Name    string
	Hash    string
	IsVideo bool
	Object  []*VideoObject `xorm:"json" json:"object,omitempty"` //视频信息
}

func init() {
	RegisterTable(Uncategorized{})
}

// AddOrUpdateVideo ...
func AddOrUpdateUncategorized(uncat *Uncategorized) (e error) {
	log.Printf("%+v", *uncat)
	if uncat.ID != "" {
		if _, err := DB().ID(uncat.ID).Update(uncat); err != nil {
			return err
		}
		return nil
	}
	if _, err := DB().InsertOne(uncat); err != nil {
		return err
	}
	return nil
}
