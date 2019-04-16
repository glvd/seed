package model

// VideoGroup ...
type VideoGroup struct {
	Index     string `json:"index"`     //索引
	Sharpness string `json:"sharpness"` //清晰度
	//Type      string         `json:"type"`                         //类型：film，FanDrama
	Output       string `json:"output"`        //输出：3D，2D
	Season       string `json:"season"`        //季
	TotalEpisode string `json:"total_episode"` //总集数
	Episode      string `json:"episode"`       //集数
	//VR        string         `xorm:"vr" json:"vr"`                 //VR格式：Half-SBS：左右半宽,Half-OU：上下半高,SBS：左右全宽
	Language string `json:"language"` //语言
	Caption  string `json:"caption"`  //字幕
	//Group    string         `json:"group"`                        //分组
	Sliced bool           `json:"sliced"`                       //切片
	HLS    *HLS           `xorm:"json" json:"hls,omitempty"`    //切片信息
	Object []*VideoObject `xorm:"json" json:"object,omitempty"` //视频信息
}

// HLS ...
type HLS struct {
	Encrypt     bool   `json:"encrypt"`      //加密
	Key         string `json:"key"`          //秘钥
	M3U8        string `json:"m3u8"`         //M3U8名
	SegmentFile string `json:"segment_file"` //ts切片名
}

//// NewVideoGroup ...
//func NewVideoGroup() *VideoGroup {
//	return &VideoGroup{
//		Sharpness: "",
//		Sliced:    false,
//		HLS:       nil,
//		Object:    nil,
//	}
//}
