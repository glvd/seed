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

	Bangumi        string        `json:"bangumi"`                    //番組
	Intro          string        `xorm:"varchar(2048)" json:"intro"` //简介
	Alias          []string      `xorm:"json" json:"alias"`          //别名，片名
	Role           []string      `xorm:"json" json:"role"`           //主演
	VideoGroupList []*ObjectList `json:"video_group_list"`           //视频
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

// LoadOld ...
func LoadOld(engine *xorm.Engine) []*Video {
	var e error
	var videos = new([]*Video)
	e = engine.Find(videos)
	if e != nil {
		log.Error(e)
		return nil
	}
	return *videos
	//ret := make(map[string]*Object)
	//for _, v := range *tables {
	//	if len(v.VideoGroupList) > 0 {
	//		if len(v.VideoGroupList[0].Object) > 0 {
	//			var obj Object
	//			obj = v.VideoGroupList[0].Object[0]
	//			ret[v.Bangumi] = &obj
	//		}
	//	}
	//}
	//
	//return ret
}

func GetObject(video *Video) *Object {
	if len(video.VideoGroupList) > 0 {
		if len(video.VideoGroupList[0].Object) > 0 {
			var obj Object
			obj = video.VideoGroupList[0].Object[0]
			return &obj
		}
	}
	return nil
}
