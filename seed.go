package seed

type Extend struct {
	Path    string `json:"path"`
	Message string `json:"message"`
}

type VideoSource struct {
	Bangumi    string    `json:"bangumi"`               //番号
	FilePath   []string  `json:"file_path"`             //存放路径
	Slice      bool      `json:"slice"`                 //是否切片
	PosterPath string    `json:"poster_path,omitempty"` //海报路径
	ExtendList []*Extend `json:"extend_list,omitempty"` //扩展信息
	Role       []string  `json:"role,omitempty"`        //角色列表
	Sharpness  string    `json:"sharpness,omitempty"`   //清晰度
	Publish    string    `json:"publish,omitempty"`     //发布日期
} //上传视频JSON配置

type VideoLink struct {
	Hash string `json:"hash"`
	Name string `json:"name"`
	Size uint64 `json:"size"`
	Type int    `json:"type"`
} //视频IPFS地址信息

type VideoGroup struct {
	Sharpness string       `json:"sharpness"`  //清晰度
	Sliced    bool         `json:"sliced"`     //切片
	VideoLink VideoLink    `json:"video_link"` //视频源
	PlayList  []*VideoLink `json:"play_list"`  //具体信息
} //整套片源

type Video struct {
	VideoInfo      VideoInfo     `json:"video_info"`       //基本信息
	VideoGroupList []*VideoGroup `json:"video_group_list"` //多套片源
}

type VideoInfo struct {
	Bangumi string   `json:"bangumi"` //番号
	Poster  string   `json:"poster"`  //海报
	Role    []string `json:"role"`    //主演
	Publish string   `json:"publish"` //发布日期
} //视频信息
