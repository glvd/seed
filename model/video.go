package model

import "golang.org/x/xerrors"

// Video ...
type Video struct {
	Model          `xorm:"extends"`
	*VideoInfo     `xorm:"extends"`
	VideoGroupList []*VideoGroup `json:"video_group_list"`
	SourceInfoList []*SourceInfo `json:"source_info_list"`
	Peers          []*SourcePeer `xorm:"-" json:"peers"`
}

// VideoInfo ...
type VideoInfo struct {
	Bangumi      string   `xorm:"unique index 'bangumi'" json:"bangumi"` //番組
	Type         string   `json:"type"`                                  //类型：film，FanDrama
	Output       string   `json:"output"`                                //输出：3D，2D
	VR           string   `xorm:"vr" json:"vr"`                          //VR格式：左右，上下，平面
	Thumb        string   `json:"thumb"`                                 //缩略图
	Intro        string   `json:"intro"`                                 //简介
	Alias        []string `xorm:"json" json:"alias"`                     //别名，片名
	Language     string   `json:"language"`                              //语言
	Caption      string   `json:"caption"`                               //字幕
	Poster       string   `json:"poster"`                                //海报
	Role         []string `xorm:"json" json:"role"`                      //主演
	Director     string   `json:"director"`                              //导演
	Season       string   `json:"season"`                                //季
	Episode      string   `bson:"episode"`                               //集数
	TotalEpisode string   `bson:"total_episode"`                         //总集数
	Sharpness    string   `json:"sharpness"`                             //清晰度
	Group        string   `json:"group"`                                 //分组
	Publish      string   `json:"publish"`                               //发布日期
}

// AddPeers ...
func (v *Video) AddPeers(p ...*SourcePeer) {
	v.Peers = append(v.Peers, p...)
}

// AddSourceInfo ...
func (v *Video) AddSourceInfo(info *SourceInfoDetail) {
	addSourceInfo(v, info)
}

// FindVideo ...
func FindVideo(ban string, video *Video) (e error) {
	if b, err := db.Where("bangumi = ?", ban).Get(video); err != nil || !b {
		return xerrors.New("video not found")
	}
	return nil
}

// AddVideo ...
func AddVideo(video *Video) (e error) {
	if _, err := db.InsertOne(video); err != nil {
		return err
	}
	return nil
}
