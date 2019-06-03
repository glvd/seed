package model

// VideoGroup ...
type VideoGroup struct {
	Model       `xorm:"extends"`
	Index       string         `json:"index"`                        //索引
	Checksum    string         `json:"checksum"`                     //文件Hash
	Sliced      bool           `json:"sliced"`                       //切片
	Encrypt     bool           `json:"encrypt"`                      //加密
	Key         string         `json:"key"`                          //秘钥
	M3U8        string         `json:"m3u8"`                         //M3U8名
	SegmentFile string         `json:"segment_file"`                 //ts切片名
	Object      []*VideoObject `xorm:"json" json:"object,omitempty"` //视频信息
}

// HLS ...
type HLS struct {
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
