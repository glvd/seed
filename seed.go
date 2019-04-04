package seed

import (
	"crypto/sha1"
	"fmt"
	"github.com/godcong/go-ipfs-restapi"
	"github.com/json-iterator/go"
)

// Extend ...
type Extend struct {
	Path    string `json:"path"`
	Message string `json:"message"`
}

// HLS ...
type HLS struct {
	Encrypt     bool   `json:"encrypt"`      //加密
	Key         string `json:"key"`          //秘钥
	M3U8        string `json:"m3u8"`         //M3U8名
	SegmentFile string `json:"segment_file"` //ts切片名
}

// VideoSource ...
type VideoSource struct {
	Bangumi      string    `json:"bangumi"`       //番号
	Type         string    `json:"type"`          //类型：film，FanDrama
	Output       string    `json:"output"`        //输出：3D，2D
	VR           string    `json:"vr"`            //VR格式：左右，上下，平面
	Thumb        string    `json:"thumb"`         //缩略图
	Intro        string    `json:"intro"`         //简介
	Alias        []string  `json:"alias"`         //别名，片名
	VideoEncode  string    `json:"video_encode"`  //视频编码
	AudioEncode  string    `json:"audio_encode"`  //音频编码
	Files        []string  `json:"files"`         //存放路径
	Slice        bool      `json:"slice"`         //是否HLS切片
	HLS          *HLS      `json:"hls"`           //HLS信息
	PosterPath   string    `json:"poster_path"`   //海报路径
	ExtendList   []*Extend `json:"extend_list"`   //扩展信息
	Role         []string  `json:"role"`          //角色列表
	Director     string    `json:"director"`      //导演
	Season       string    `json:"season"`        //季
	Episode      string    `bson:"episode"`       //集数
	TotalEpisode string    `bson:"total_episode"` //总集数
	Sharpness    string    `json:"sharpness"`     //清晰度
	//Group        string    `json:"group"`         //分组,废弃
	Publish string `json:"publish"` //发布日期
}

// VideoLink ...
type VideoLink struct {
	Hash string `json:"hash,omitempty"`
	Name string `json:"name,omitempty"`
	Size uint64 `json:"size,omitempty"`
	Type int    `json:"type,omitempty"`
}

// VideoGroup ...
type VideoGroup struct {
	Index     string    `json:"index"`            //索引
	Sharpness string    `json:"sharpness"`        //清晰度
	Sliced    bool      `json:"sliced"`           //切片
	HLS       *HLS      `json:"hls,omitempty"`    //切片信息
	Object    []*Object `json:"object,omitempty"` //视频信息
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

// SourceInfo ...
type SourceInfo struct {
	ID              string   `json:"id"`
	PublicKey       string   `json:"public_key"`
	Addresses       []string `json:"addresses"` //一组节点源列表
	AgentVersion    string   `json:"agent_version"`
	ProtocolVersion string   `json:"protocol_version"`
}

// Video ...
type Video struct {
	*VideoInfo     `json:",inline"` //基本信息
	VideoGroupList []*VideoGroup    `json:"video_group_list"` //多套片源
	SourceInfoList []*SourceInfo    `json:"source_info_list"` //节点源数据
}

// Link ...
type Link struct {
	Hash string `json:"hash"`
	Name string `json:"name"`
	Size uint64 `json:"size"`
	Type int    `json:"type"`
}

// Object ...
type Object struct {
	Links []*Link `json:"links,omitempty"`
	Link  `json:",inline"`
}

// VideoList ...
var VideoList = LoadVideo()

// LoadVideo ...
func LoadVideo() []*Video {
	var videos []*Video
	e := ReadJSON("video.json", &videos)
	if e != nil {
		return nil
	}
	return videos
}

// ListVideoGet ...
func ListVideoGet(source *VideoSource) *Video {
	for _, v := range VideoList {
		if v.Bangumi == source.Bangumi {
			return v
		}
	}
	return NewVideo(source)
}

// VideoListAdd ...
func VideoListAdd(source *VideoSource, video *Video) {
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

// NewVideo ...
func NewVideo(source *VideoSource) *Video {
	alias := []string{}
	if source.Alias != nil {
		alias = source.Alias
	}
	return &Video{
		VideoInfo: &VideoInfo{
			Bangumi: source.Bangumi,
			Alias:   alias,
			Role:    source.Role,
			Publish: source.Publish,
		},
		VideoGroupList: nil,
		SourceInfoList: nil,
	}
}

// AddSourceInfo ...
func AddSourceInfo(video *Video, info *SourceInfo) {
	if video.SourceInfoList == nil {
		video.SourceInfoList = []*SourceInfo{info}
		return
	}
	for idx, value := range video.SourceInfoList {
		if value.ID == info.ID {
			video.SourceInfoList[idx] = info
			return
		}
	}
	video.SourceInfoList = append(video.SourceInfoList, info)
}

// NewVideoGroup ...
func NewVideoGroup() *VideoGroup {
	return &VideoGroup{
		Sharpness: "",
		Sliced:    false,
		HLS:       nil,
		Object:    nil,
	}
}

// LinkObjectToObject ...
func LinkObjectToObject(obj interface{}) *Object {
	if v, b := obj.(*Object); b {
		return v
	}
	return &Object{}
}

// AddRetToLink ...
func AddRetToLink(obj *Object, ret *api.AddRet) *Object {
	if obj != nil {
		obj.Link.Hash = ret.Hash
		obj.Link.Name = ret.Name
		obj.Link.Size = ret.Size
		obj.Link.Type = 2
		return obj
	}
	return &Object{
		Link: Link{
			Hash: ret.Hash,
			Name: ret.Name,
			Size: ret.Size,
			Type: 2,
		},
	}
}

// AddRetToLinks ...
func AddRetToLinks(obj *Object, ret *api.AddRet) *Object {
	if obj != nil {
		obj.Links = append(obj.Links, &Link{
			Hash: ret.Hash,
			Name: ret.Name,
			Size: ret.Size,
			Type: 2,
		})
		return obj
	}
	return &Object{
		Links: []*Link{
			{
				Hash: ret.Hash,
				Name: ret.Name,
				Size: ret.Size,
				Type: 2,
			},
		},
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
