package seed

type VideoSource struct {
	Bangumi   string        `json:"bangumi"`   //番号
	Path      string        `json:"path"`      //存放路径
	Poster    string        `json:"poster"`    //海报
	Role      []interface{} `json:"role"`      //主演
	Sharpness string        `json:"sharpness"` //清晰度
	Publish   string        `json:"publish"`   //发布日期
}

type VideoLink struct {
	Hash string `json:"hash"`
	Name string `json:"name"`
	Size uint64 `json:"size"`
	Type int    `json:"type"`
}

type VideoLinkInfo struct {
	Sharpness string    `json:"sharpness"` //清晰度
	Sliced    bool      `json:"sliced"`    //切片
	Type      string    `json:"type"`      //目录，文件
	Index     string    `json:"index"`     //文件
	VideoLink VideoLink `json:"video_link_list"`
}

type Video struct {
	Bangumi string        `json:"bangumi"` //番号
	Poster  string        `json:"poster"`  //海报
	Role    []interface{} `json:"role"`    //主演
	Publish string        `json:"publish"` //发布日期
}

type VideoInfo struct {
}
