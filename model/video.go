package model

// Video ...
type Video struct {
	Model          `xorm:"extends"`
	*VideoInfo     `xorm:"extends"`
	VideoGroupList []*VideoGroup `json:"video_group_list"`
	SourceInfoList []*SourceInfo `json:"source_info_list"`
	Peers          []string      `xorm:"-" json:"peers"`
}

// VideoInfo ...
type VideoInfo struct {
	Bangumi      string   `json:"bangumi"`           //番組
	Type         string   `json:"type"`              //类型：film，FanDrama
	Output       string   `json:"output"`            //输出：3D，2D
	VR           string   `xorm:"vr" json:"vr"`      //VR格式：左右，上下，平面
	Thumb        string   `json:"thumb"`             //缩略图
	Intro        string   `json:"intro"`             //简介
	Alias        []string `xorm:"json" json:"alias"` //别名，片名
	Language     string   `json:"language"`          //语言
	Caption      string   `json:"caption"`           //字幕
	Poster       string   `json:"poster"`            //海报
	Role         []string `xorm:"json" json:"role"`  //主演
	Director     string   `json:"director"`          //导演
	Season       string   `json:"season"`            //季
	Episode      string   `bson:"episode"`           //集数
	TotalEpisode string   `bson:"total_episode"`     //总集数
	Sharpness    string   `json:"sharpness"`         //清晰度
	Group        string   `json:"group"`             //分组
	Publish      string   `json:"publish"`           //发布日期
}
