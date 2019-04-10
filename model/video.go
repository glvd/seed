package model

// Video ...
type Video struct {
	Model     `xorm:",extends"`
	VideoInfo `xorm:",extends"`
}

// VideoInfo ...
type VideoInfo struct {
	Bangumi      string   `json:"bangumi"`       //番組
	Type         string   `json:"type"`          //类型：film，FanDrama
	Output       string   `json:"output"`        //输出：3D，2D
	VR           string   `json:"vr"`            //VR格式：左右，上下，平面
	Thumb        string   `json:"thumb"`         //缩略图
	Intro        string   `json:"intro"`         //简介
	Alias        []string `json:"alias"`         //别名，片名
	Language     string   `json:"language"`      //语言
	Caption      string   `json:"caption"`       //字幕
	Poster       string   `json:"poster"`        //海报
	Role         []string `json:"role"`          //主演
	Director     string   `json:"director"`      //导演
	Season       string   `json:"season"`        //季
	Episode      string   `bson:"episode"`       //集数
	TotalEpisode string   `bson:"total_episode"` //总集数
	Sharpness    string   `json:"sharpness"`     //清晰度
	Group        string   `json:"group"`         //分组
	Publish      string   `json:"publish"`       //发布日期
}
