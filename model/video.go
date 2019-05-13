package model

import (
	log "github.com/sirupsen/logrus"
)

// Video ...
type Video struct {
	Model          `xorm:"extends"`
	Sync           bool
	*VideoBase     `xorm:"extends"`
	VideoGroupList []*VideoGroup `xorm:"json" json:"video_group_list"`
	SourceInfoList []*SourceInfo `xorm:"json" json:"source_info_list"`
	SourcePeerList []*SourcePeer `xorm:"json" json:"source_peer_list"`
}

// VideoBase ...
type VideoBase struct {
	Bangumi      string   `xorm:"unique index bangumi" json:"bangumi"` //番組
	Thumb        string   `json:"thumb"`                               //缩略图
	Intro        string   `xorm:"varchar(2048)" json:"intro"`          //简介
	Alias        []string `xorm:"json" json:"alias"`                   //别名，片名
	SourceHash   string   `json:"source_hash"`                         //原片地址
	Poster       string   `json:"poster"`                              //海报
	Role         []string `xorm:"json" json:"role"`                    //主演
	Director     []string `xorm:"json" json:"director"`                //导演
	Season       string   `json:"season,omitempty"`                    //季
	TotalEpisode string   `json:"total_episode,omitempty"`             //总集数
	Episode      string   `json:"episode,omitempty"`                   //集数
	Publish      string   `json:"publish"`                             //发布日期
}

// VideoInfo ...
type VideoInfo struct {
	Type     string `json:"type"`         //类型：film，FanDrama
	Output   string `json:"output"`       //输出：3D，2D
	VR       string `xorm:"vr" json:"vr"` //VR格式：Half-SBS：左右半宽,Half-OU：上下半高,SBS：左右全宽
	Language string `json:"language"`     //语言
	Caption  string `json:"caption"`      //字幕
	Group    string `json:"group"`        //分组
}

func init() {
	RegisterTable(Video{})
}

// AddPeers ...
func (v *Video) AddPeers(p ...*SourcePeerDetail) {
	for _, value := range p {
		v.SourcePeerList = append(v.SourcePeerList, &SourcePeer{SourcePeerDetail: value})
	}
}

// AddSourceInfo ...
func (v *Video) AddSourceInfo(info *SourceInfoDetail) {
	addSourceInfo(v, info)
}

// FindVideo ...
func FindVideo(ban string, video *Video, check bool) (b bool, e error) {
	if check {
		return DB().Where("sync = ?", check).Where("bangumi like ?", "%"+ban+"%").Get(video)
	}
	return DB().Where("bangumi like ?", "%"+ban+"%").Get(video)
}

// Top ...
func Top(video *Video) (b bool, e error) {
	return DB().OrderBy("created_at desc").Get(video)
}

// AllVideos ...
func AllVideos(check bool) (v []*Video, e error) {
	var videos = new([]*Video)
	if check {
		if e = DB().Where("sync = ?", check).Find(videos); e != nil {
			return
		}
	} else {
		if e = DB().Find(videos); e != nil {
			return
		}
	}
	v = *videos
	return
}

// DeepFind ...
func DeepFind(s string, video *Video) (b bool, e error) {
	b, e = DB().Where("bangumi = ?", s).Get(video)
	if e != nil || !b {
		like := "%" + s + "%"
		return DB().Where("bangumi like ? ", like).
			Or("alias like ?", like).
			Or("role like ?", like).
			Get(video)
	}
	return b, e
}

// AddOrUpdateVideo ...
func AddOrUpdateVideo(video *Video) (e error) {
	log.Printf("%+v", *video)
	if video.ID != "" {
		log.Debug("update")
		if _, err := DB().ID(video.ID).Update(video); err != nil {
			return err
		}
		return nil
	}
	if _, err := DB().InsertOne(video); err != nil {
		return err
	}
	return nil
}
