package seed

import (
	"crypto/sha1"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/yinhevr/seed/model"
	"strings"
)

// Extend ...
type Extend struct {
	Path    string `json:"path"`
	Message string `json:"message"`
}

// VideoSource ...
type VideoSource struct {
	Bangumi      string     `json:"bangumi"`       //番号
	SourceHash   string     `json:"source_hash"`   //原片hash
	Type         string     `json:"type"`          //类型：film，FanDrama
	Output       string     `json:"output"`        //输出：3D，2D
	VR           string     `json:"vr"`            //VR格式：左右，上下，平面
	Thumb        string     `json:"thumb"`         //缩略图
	Intro        string     `json:"intro"`         //简介
	Alias        []string   `json:"alias"`         //别名，片名
	VideoEncode  string     `json:"video_encode"`  //视频编码
	AudioEncode  string     `json:"audio_encode"`  //音频编码
	Files        []string   `json:"files"`         //存放路径
	HashFiles    []string   `json:"hash_files"`    //已上传Hash
	CheckFiles   []string   `json:"check_files"`   //uncategorized checksum
	Slice        bool       `json:"slice"`         //是否HLS切片
	HLS          *model.HLS `json:"hls"`           //HLS信息
	PosterPath   string     `json:"poster_path"`   //海报路径
	Poster       string     `json:"poster"`        //海报HASH
	ExtendList   []*Extend  `json:"extend_list"`   //扩展信息
	Role         []string   `json:"role"`          //角色列表
	Director     []string   `json:"director"`      //导演
	Systematics  string     `json:"systematics"`   //分级
	Season       string     `json:"season"`        //季
	Episode      string     `json:"episode"`       //集数
	TotalEpisode string     `json:"total_episode"` //总集数
	Sharpness    string     `json:"sharpness"`     //清晰度
	Publish      string     `json:"publish"`       //发行日
	Language     string     `json:"language"`      //语言
	Caption      string     `json:"caption"`       //字幕
}

// VideoLink ...
type VideoLink struct {
	Hash string `json:"hash,omitempty"`
	Name string `json:"name,omitempty"`
	Size uint64 `json:"size,omitempty"`
	Type int    `json:"type,omitempty"`
}

// VideoList ...
var VideoList = LoadVideo()

// LoadVideo ...
func LoadVideo() []*model.Video {
	var videos []*model.Video
	e := ReadJSON("video.json", &videos)
	if e != nil {
		return nil
	}
	return videos
}

// ListVideoGet ...
func ListVideoGet(source *VideoSource) *model.Video {
	for _, v := range VideoList {
		if v.Bangumi == source.Bangumi {
			return v
		}
	}
	return NewVideo(source)
}

// VideoListAdd ...
func VideoListAdd(source *VideoSource, video *model.Video) {
	for i, v := range VideoList {
		if v.Bangumi == source.Bangumi {
			VideoList[i] = video
			return
		}
	}
	VideoList = append(VideoList, video)
}

// SaveVideos ...
func SaveVideos() (e error) {
	e = WriteJSON("video.json", VideoList)
	if e != nil {
		return e
	}
	return nil
}

func parseVideoBase(video *model.Video, source *VideoSource) {
	if video == nil {
		return
	}
	alias := []string{}
	aliasS := ""
	if source.Alias != nil || len(source.Alias) > 0 {
		alias = source.Alias
		aliasS = alias[0]
	}

	role := []string{}
	roleS := ""
	if source.Role != nil || len(source.Role) > 0 {
		role = source.Role
		roleS = role[0]
	}

	director := []string{}
	if source.Director != nil {
		director = source.Director
	}

	intro := source.Intro
	if intro != "" {
		intro = aliasS + " " + roleS
	}

	video.VideoBase = &model.VideoBase{
		Bangumi: source.Bangumi,
		//Type:         source.Type,
		//Output:       source.Output,
		//VR:           source.VR,
		Thumb:    source.Thumb,
		Intro:    intro,
		Alias:    alias,
		Role:     role,
		Director: director,
		//Language:     source.Language,
		//Caption:      source.Caption,
		SourceHash: source.SourceHash,

		Season:       source.Season,
		Episode:      source.Episode,
		TotalEpisode: source.TotalEpisode,
		Publish:      source.Publish,
	}
}

// NewVideo ...
func NewVideo(source *VideoSource) *model.Video {
	alias := []string{}
	if source.Alias != nil {
		alias = source.Alias
	}
	return &model.Video{
		VideoBase: &model.VideoBase{
			Bangumi: strings.ToUpper(source.Bangumi),
			//Type:         source.Type,
			//Output:       source.Output,
			//VR:           source.VR,
			Thumb: source.Thumb,
			Intro: source.Intro,
			Alias: alias,
			//Language:     source.Language,
			//Caption:      source.Caption,
			Role:         source.Role,
			Director:     source.Director,
			Season:       source.Season,
			Episode:      source.Episode,
			TotalEpisode: source.TotalEpisode,
			Publish:      source.Publish,
		},
		VideoGroupList: nil,
		SourceInfoList: nil,
	}
}

// Hash ...
func Hash(v interface{}) string {
	bytes, e := jsoniter.Marshal(v)
	if e != nil {
		return ""
	}
	return fmt.Sprintf("%x", sha1.Sum([]byte(bytes)))
}
