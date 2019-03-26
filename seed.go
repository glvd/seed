package seed

import "github.com/ipfs/go-ipfs-api"

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
	Bangumi    string    `json:"bangumi"`               //番号
	FilePath   []string  `json:"file_path"`             //存放路径
	Slice      bool      `json:"slice"`                 //是否HLS切片
	HLS        HLS       `json:"hls,omitempty"`         //HLS信息
	PosterPath string    `json:"poster_path,omitempty"` //海报路径
	ExtendList []*Extend `json:"extend_list,omitempty"` //扩展信息
	Role       []string  `json:"role,omitempty"`        //角色列表
	Sharpness  string    `json:"sharpness,omitempty"`   //清晰度
	Group      string    `json:"group"`                 //分组
	Publish    string    `json:"publish,omitempty"`     //发布日期
} //上传视频JSON配置

type VideoLink struct {
	Hash string `json:"hash"`
	Name string `json:"name"`
	Size uint64 `json:"size"`
	Type int    `json:"type"`
} //视频IPFS地址信息

type VideoGroup struct {
	Sharpness string       `json:"sharpness,omitempty"` //清晰度
	Sliced    bool         `json:"sliced,omitempty"`    //切片
	HLS       HLS          `json:"hls,omitempty"`
	VideoLink VideoLink    `json:"video_link,omitempty"` //视频源
	PlayList  []*VideoLink `json:"play_list,omitempty"`  //具体信息
} //整套片源

type VideoInfo struct {
	Bangumi string   `json:"bangumi"` //番号
	Poster  string   `json:"poster"`  //海报
	Role    []string `json:"role"`    //主演
	Publish string   `json:"publish"` //发布日期
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
		panic(e)
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
	return &Video{
		VideoInfo: &VideoInfo{
			Bangumi: source.Bangumi,
			//Poster:  source.PosterPath,
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
		VideoLink: VideoLink{},
		PlayList:  nil,
	}
}

func LsLinkToVideoLink(link *shell.LsLink) *VideoLink {
	return &VideoLink{
		Hash: link.Hash,
		Name: link.Name,
		Size: link.Size,
		Type: link.Type,
	}
}
