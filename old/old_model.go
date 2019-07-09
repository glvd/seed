package old

import (
	"github.com/go-xorm/xorm"
	"github.com/godcong/go-trait"
	"github.com/yinhevr/seed/model"
)

var log = trait.NewZapSugar()

// Video ...
type Video struct {
	//FindNo         string        `json:"find_no"`          //查找号
	Bangumi        string        `json:"bangumi"`          //番組
	VideoGroupList []*ObjectList `json:"video_group_list"` //视频
}

// ObjectList ...
type ObjectList struct {
	Index        string   `json:"index"`
	Sharpness    string   `json:"sharpness"`
	Output       string   `json:"output"`
	Season       string   `json:"season"`
	TotalEpisode string   `json:"total_episode"`
	Episode      string   `json:"episode"`
	Language     string   `json:"language"`
	Caption      string   `json:"caption"`
	Sliced       bool     `json:"sliced"`
	Object       []Object `json:"object"`
}

// Object ...
type Object struct {
	Links []*model.VideoLink `json:"links,omitempty"`
	Link  *model.VideoLink   `xorm:"extends"  json:",inline,omitempty"`
}

// LoadFrom ...
func LoadFrom(path string) map[string]*Object {
	engine, e := xorm.NewEngine("sqlite3", path)
	if e != nil {
		log.Error(e)
		return nil
	}
	e = engine.Sync2(Video{})
	if e != nil {
		log.Error(e)
		return nil
	}
	var tables = new([]*Video)
	e = engine.Find(tables)
	if e != nil {
		log.Error(e)
		return nil
	}

	ret := make(map[string]*Object)
	for _, v := range *tables {
		if len(v.VideoGroupList) > 0 {
			if len(v.VideoGroupList[0].Object) > 0 {
				var obj Object
				obj = v.VideoGroupList[0].Object[0]
				ret[v.Bangumi] = &obj
			}
		}
	}

	return ret
}
