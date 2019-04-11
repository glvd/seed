package model

import (
	log "github.com/sirupsen/logrus"
)

// Video ...
type Video struct {
	Model          `xorm:"extends"`
	*VideoInfo     `xorm:"extends"`
	VideoGroupList []*VideoGroup `xorm:"json" json:"video_group_list"`
	SourceInfoList []*SourceInfo `xorm:"json" json:"source_info_list"`
	SourcePeerList []*SourcePeer `xorm:"json" json:"source_peer_list"`
}

// VideoInfo ...
type VideoInfo struct {
	Bangumi      string   `xorm:"unique index 'bangumi'" json:"bangumi"` //番組
	Type         string   `json:"type"`                                  //类型：film，FanDrama
	Output       string   `json:"output"`                                //输出：3D，2D
	VR           string   `xorm:"vr" json:"vr"`                          //VR格式：Half-SBS：左右半宽,Half-OU：上下半高,SBS：左右全宽
	Thumb        string   `json:"thumb"`                                 //缩略图
	Intro        string   `json:"intro"`                                 //简介
	Alias        []string `xorm:"json" json:"alias"`                     //别名，片名
	Language     string   `json:"language"`                              //语言
	Caption      string   `json:"caption"`                               //字幕
	Poster       string   `json:"poster"`                                //海报
	Role         []string `xorm:"json" json:"role"`                      //主演
	Director     string   `json:"director"`                              //导演
	Season       string   `json:"season"`                                //季
	Episode      string   `json:"episode"`                               //集数
	TotalEpisode string   `json:"total_episode"`                         //总集数
	Group        string   `json:"group"`                                 //分组
	Publish      string   `json:"publish"`                               //发布日期
}

func init() {
	RegisterTable(&Video{})
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
func FindVideo(ban string, video *Video) (b bool, e error) {
	return DB().Where("bangumi = ?", ban).Get(video)
}

// AddVideo ...
func AddVideo(video *Video) (e error) {
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
