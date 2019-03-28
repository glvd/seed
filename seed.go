package seed

import (
	"github.com/godcong/go-ipfs-restapi"
)

type Extend struct {
	Path    string `json:"path"`
	Message string `json:"message"`
}

type HLS struct {
	Encrypt    bool   `json:"encrypt"`     //加密
	M3U8       string `json:"m3u8"`        //M3U8名
	OutputName string `json:"output_name"` //ts名
}

type VideoSource struct {
	Bangumi     string    `json:"bangumi"`               //番号
	Alias       []string  `json:"alias"`                 //别名，片名
	VideoEncode string    `json:"video_encode"`          //视频编码
	AudioEncode string    `json:"audio_encode"`          //音频编码
	Files       []string  `json:"files"`                 //存放路径
	Slice       bool      `json:"slice"`                 //是否HLS切片
	HLS         HLS       `json:"hls,omitempty"`         //HLS信息
	PosterPath  string    `json:"poster_path,omitempty"` //海报路径
	ExtendList  []*Extend `json:"extend_list,omitempty"` //扩展信息
	Role        []string  `json:"role,omitempty"`        //角色列表
	Sharpness   string    `json:"sharpness,omitempty"`   //清晰度
	Group       string    `json:"group"`                 //分组
	Publish     string    `json:"publish,omitempty"`     //发布日期
} //上传视频JSON配置

type VideoLink struct {
	Hash string `json:"hash,omitempty"`
	Name string `json:"name,omitempty"`
	Size uint64 `json:"size,omitempty"`
	Type int    `json:"type,omitempty"`
} //视频IPFS地址信息

type VideoGroup struct {
	Sharpness string          `json:"sharpness,omitempty"` //清晰度
	Sliced    bool            `json:"sliced,omitempty"`    //切片
	HLS       HLS             `json:"hls,omitempty"`
	Object    *api.ListObject `json:"object,omitempty"` //视频信息
	//PlayList  []*VideoLink `json:"play_list,omitempty"`  //具体信息
} //整套片源

type VideoInfo struct {
	Bangumi  string   `json:"bangumi"`  //番号
	Alias    []string `json:"alias"`    //别名，片名
	Language string   `json:"language"` //语言
	Caption  string   `json:"caption"`  //字幕
	Poster   string   `json:"poster"`   //海报
	Role     []string `json:"role"`     //主演
	Publish  string   `json:"publish"`  //发布日期
} //视频信息

type Video struct {
	*VideoInfo     `json:",inline"`       //基本信息
	VideoGroupList map[string]*VideoGroup `json:"video_group_list"` //多套片源
}

var VideoList = LoadVideo()

func LoadVideo() map[string]*Video {
	videos := make(map[string]*Video)
	e := ReadJSON("video.json", &videos)
	if e != nil {
		return make(map[string]*Video)
	}
	return videos
}

func GetVideo(source *VideoSource) *Video {
	if video, b := VideoList[source.Bangumi]; b {
		return video
	}
	return NewVideo(source)
}

func SetVideo(source *VideoSource, video *Video) {
	VideoList[source.Bangumi] = video
}

func SaveVideos() (e error) {
	e = WriteJSON("video.json", VideoList)
	if e != nil {
		return e
	}
	return nil
}

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
	}
}
func NewVideoGroup() *VideoGroup {
	return &VideoGroup{
		Sharpness: "",
		Sliced:    false,
		HLS:       HLS{},
		Object:    nil,
	}
}

func AddRetToObject(ret *api.AddRet) *api.ListObject {
	return &api.ListObject{
		ListLink: api.ListLink{
			Hash: ret.Hash,
			Name: ret.Name,
			Size: ret.Size,
			Type: 2,
		},
	}
}
