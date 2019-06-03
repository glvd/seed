package model

// Video ...
type Video struct {
	Model        `xorm:"extends"`
	FindNo       string   `json:"find_no"`                    //查找号
	Bangumi      string   `json:"bangumi"`                    //番組
	Thumb        string   `json:"thumb"`                      //缩略图
	Intro        string   `xorm:"varchar(2048)" json:"intro"` //简介
	Alias        []string `xorm:"json" json:"alias"`          //别名，片名
	SourceHash   string   `json:"source_hash"`                //原片地址
	M3U8Hash     string   `json:"m3u8_hash"`                  //切片地址
	PosterHash   string   `json:"poster_hash"`                //海报地址
	Role         []string `xorm:"json" json:"role"`           //主演
	Director     []string `xorm:"json" json:"director"`       //导演
	Season       string   `json:"season,omitempty"`           //季
	TotalEpisode string   `json:"total_episode,omitempty"`    //总集数
	Episode      string   `json:"episode,omitempty"`          //集数
	Publish      string   `json:"publish"`                    //发布日期
	Type         string   `json:"type"`                       //类型：film，FanDrama
	Format       string   `json:"format"`                     //输出格式：3D，2D
	VR           string   `xorm:"vr" json:"vr"`               //VR格式：Half-SBS：左右半宽,Half-OU：上下半高,SBS：左右全宽
	Language     string   `json:"language"`                   //语言
	Caption      string   `json:"caption"`                    //字幕
	Group        string   `json:"group"`                      //分组
	Index        string   `json:"index"`                      //索引
	Sharpness    string   `json:"sharpness"`                  //清晰度
	Sliced       bool     `json:"sliced"`                     //切片
	HLS          HLS      `xorm:"json" json:"hls,omitempty"`  //切片信息
	Sync         bool     `xorm:"notnull default(0)"`         //是否已同步
	Visit        uint64   `xorm:"notnull default(0)"`         //访问统计
	//VideoGroupList []*VideoGroup `xorm:"json" json:"video_group_list"`
	//SourceInfoList []*SourceInfo `xorm:"json" json:"source_info_list"`
	//SourcePeerList []*SourcePeer `xorm:"json" json:"source_peer_list"`
}

func init() {
	RegisterTable(Video{})
}

// AddPeers ...
func (v *Video) AddPeers(p ...*SourcePeerDetail) {
	for _, value := range p {
		v.SourcePeerList = append(v.SourcePeerList, &SourcePeer{SourcePeerDetail: value})
	}
}

// AddSourceInfo ...
func (v *Video) AddSourceInfo(info *SourceInfoDetail) {
	addSourceInfo(v, info)
}

// FindVideo ...
func FindVideo(ban string, video *Video, check bool) (b bool, e error) {
	if check {
		return DB().Where("sync = ?", !check).Where("bangumi like ?", "%"+ban+"%").Get(video)
	}
	return DB().Where("bangumi like ?", "%"+ban+"%").Get(video)
}

// Top ...
func Top(video *Video) (b bool, e error) {
	return DB().OrderBy("created_at desc").Get(video)
}

// AllVideos ...
func AllVideos(check bool) (v []*Video, e error) {
	var videos = new([]*Video)
	if check {
		if e = DB().Where("sync = ?", !check).Find(videos); e != nil {
			return
		}
	} else {
		if e = DB().Find(videos); e != nil {
			return
		}
	}
	v = *videos
	return
}

// DeepFind ...
func DeepFind(s string, video *Video) (b bool, e error) {
	b, e = DB().Where("bangumi = ?", s).Get(video)
	if e != nil || !b {
		like := "%" + s + "%"
		return DB().Where("bangumi like ? ", like).
			Or("alias like ?", like).
			Or("role like ?", like).
			Get(video)
	}
	return b, e
}

// AddOrUpdateVideo ...
func AddOrUpdateVideo(video *Video) (e error) {
	log.Infof("%+v", *video)
	if video.ID != "" {
		log.Debug("update")
		if _, err := DB().ID(video.ID).Update(video); err != nil {
			return err
		}
		return nil
	}
	if _, err := DB().InsertOne(video); err != nil {
		return err
	}
	return nil
}

// Visited ...
func Visited(video *Video) (err error) {
	video.Visit++
	if _, err := DB().ID(video.ID).Cols("visit").Update(video); err != nil {
		return err
	}
	return nil
}
