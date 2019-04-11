package model

// VideoGroup ...
type VideoGroup struct {
	//Model     `xorm:"extends"`
	Index     string         `json:"index"`                        //索引
	Sharpness string         `json:"sharpness"`                    //清晰度
	Sliced    bool           `json:"sliced"`                       //切片
	HLS       *HLS           `xorm:"json" json:"hls,omitempty"`    //切片信息
	Object    []*VideoObject `xorm:"json" json:"object,omitempty"` //视频信息
}

// HLS ...
type HLS struct {
	Encrypt     bool   `json:"encrypt"`      //加密
	Key         string `json:"key"`          //秘钥
	M3U8        string `json:"m3u8"`         //M3U8名
	SegmentFile string `json:"segment_file"` //ts切片名
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
