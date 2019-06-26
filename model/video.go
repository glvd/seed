package model

import (
	"github.com/go-xorm/xorm"
	"strings"
)

// Video ...
type Video struct {
	Model        `xorm:"extends"`
	FindNo       string   `json:"find_no"`                    //查找号
	Bangumi      string   `json:"bangumi"`                    //番組
	ThumbHash    string   `json:"thumb_hash"`                 //缩略图
	Intro        string   `xorm:"varchar(2048)" json:"intro"` //简介
	Alias        []string `xorm:"json" json:"alias"`          //别名，片名
	SourceHash   string   `json:"source_hash"`                //原片地址
	Key          string   `json:"key"`                        //秘钥
	M3U8         string   `xorm:"m3u8" json:"m3u8"`           //M3U8名
	M3U8Hash     string   `xorm:"m3u8_hash" json:"m3u8_hash"` //切片地址
	PosterHash   string   `json:"poster_hash"`                //海报地址
	Role         []string `xorm:"json" json:"role"`           //主演
	Director     string   `json:"director"`                   //导演
	Season       string   `json:"season,omitempty"`           //季
	TotalEpisode string   `json:"total_episode,omitempty"`    //总集数
	Episode      string   `json:"episode,omitempty"`          //集数
	Producer     string   `json:"producer"`                   //生产商
	Publisher    string   `json:"publish"`                    //发行商
	Type         string   `json:"type"`                       //类型：film，FanDrama
	Format       string   `json:"format"`                     //输出格式：3D，2D,VR(VR格式：Half-SBS：左右半宽,Half-OU：上下半高,SBS：左右全宽)
	Language     string   `json:"language"`                   //语言
	Caption      string   `json:"caption"`                    //字幕
	Group        string   `json:"group"`                      //分组
	Index        string   `json:"index"`                      //索引
	Date         string   `json:"date"`                       //发行日期
	Sharpness    string   `json:"sharpness"`                  //清晰度
	Visit        uint64   `xorm:"notnull default(0)"`         //访问数
	Series       string   `json:"series"`                     //系列
	Tags         []string `xorm:"json" json:"tags"`           //标签
	Length       string   `json:"length"`                     //时长
	//访问统计
	//HLS          HLS      `xorm:"json" json:"hls,omitempty"`  //切片信息
	//VideoGroupList []*VideoGroup `xorm:"json" json:"video_group_list"`
	//SourceInfoList []*SourceInfo `xorm:"json" json:"source_info_list"`
	//SourcePeerList []*SourcePeer `xorm:"json" json:"source_peer_list"`
}

func init() {
	RegisterTable(Video{})
}

// FindVideo ...
func FindVideo(session *xorm.Session, ban string) (video *Video, e error) {
	video = new(Video)
	ban = strings.ReplaceAll(ban, "-", "")
	ban = strings.ReplaceAll(ban, "_", "")
	ban = strings.ToUpper(ban)
	_, e = DB().Where("find_no = ?", ban).Get(video)
	return
}

// Top ...
func Top(video *Video) (b bool, e error) {
	return DB().OrderBy("created_at desc").Get(video)
}

// AllVideos ...
func AllVideos(session *xorm.Session, limit int, start ...int) (videos *[]*Video, e error) {
	videos = new([]*Video)
	session = MustSession(session)
	if limit > 0 {
		session = session.Limit(limit, start...)
	}
	if e = session.Find(videos); e != nil {
		return
	}
	return videos, nil
}

// DeepFind ...
func DeepFind(s string, video *Video) (b bool, e error) {
	s1 := strings.ReplaceAll(s, "-", "")
	s1 = strings.ReplaceAll(s1, "_", "")
	s1 = strings.ToUpper(s1)
	b, e = DB().Where("find_no = ?", s1).Get(video)
	if e != nil || !b {
		like := "%" + strings.ToUpper(s) + "%"
		return DB().Where("find_no like ? ", like).
			Or("alias like ?", like).
			Or("role like ?", like).
			Get(video)
	}
	return b, e
}

func mustStr(s *string, d string) {
	if *s == "" {
		*s = d
	}
}

// AddOrUpdateVideo ...
func AddOrUpdateVideo(video *Video) (e error) {
	var tmp Video
	found, e := DB().Where("bangumi = ?", video.Bangumi).
		Where("season = ?", video.Season).
		Where("episode = ?", video.Episode).
		Where("sharpness = ?", video.Sharpness).
		Get(&tmp)
	if e != nil {
		return e
	}
	if found {
		video.Version = tmp.Version
		video.ID = tmp.ID
		mustStr(&video.M3U8Hash, tmp.M3U8Hash)
		mustStr(&video.SourceHash, tmp.SourceHash)
		mustStr(&video.PosterHash, tmp.PosterHash)
		mustStr(&video.ThumbHash, tmp.ThumbHash)
		_, e = DB().ID(video.ID).Update(video)
		return
	}
	_, e = DB().InsertOne(video)
	return
}

// Visited ...
func Visited(video *Video) (err error) {
	video.Visit++
	if _, err := DB().ID(video.ID).Cols("visit").Update(video); err != nil {
		return err
	}
	return nil
}

// Clone ...
func (v *Video) Clone() (n *Video) {
	n = new(Video)
	*n = *v
	return
}
