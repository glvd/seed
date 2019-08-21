package seed

import (
	"bytes"
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/glvd/seed/model"
	"github.com/go-xorm/xorm"
	files "github.com/ipfs/go-ipfs-files"
	httpapi "github.com/ipfs/go-ipfs-http-client"
	"go.uber.org/atomic"
	"golang.org/x/xerrors"
)

// InfoType ...
type InfoType string

// InfoTypeNone ...
const InfoTypeNone InfoType = "none"

// InfoTypeJSON ...
const InfoTypeJSON InfoType = "json"

// InfoTypeBSON ...
const InfoTypeBSON InfoType = "bson"

// Information ...
type Information struct {
	Seeder
	InfoType     InfoType
	Path         string
	ResourcePath string
	ProcList     []string
	Start        int
	vcb          chan InformationVideoCallback
}

// Push ...
func (info *Information) Push(v interface{}) error {
	return info.pushVideoCallback(v)
}

// Option ...
func (info *Information) Option(seed Seeder) {
	informationOption(info)(seed)
}

// NewInformation ...
func NewInformation() *Information {
	info := new(Information)
	return info
}

// InformationVideoCallback ...
type InformationVideoCallback func(information *Information, v *model.Video)

// PushCallback ...
func (info *Information) pushVideoCallback(cb interface{}) error {
	if v, b := cb.(InformationVideoCallback); b {
		go func(callback InformationVideoCallback) {
			info.vcb <- callback
		}(v)
	}
	return xerrors.New("not information callback")
}

// BeforeRun ...
func (info *Information) BeforeRun(seed Seeder) {
	info.Seeder = seed
}

// AfterRun ...
func (info *Information) AfterRun(seed Seeder) {
}

func fixBson(s []byte) []byte {
	reg := regexp.MustCompile(`("_id")[ ]*[:][ ]*(ObjectId\(")[\w]{24}("\))[ ]*(,)[ ]*`)
	return reg.ReplaceAll(s, []byte(" "))
}

func video(source *VideoSource) (video *model.Video) {

	//always not null
	alias := *new([]string)
	aliasS := ""
	if source.Alias != nil && len(source.Alias) > 0 {
		alias = source.Alias
		aliasS = alias[0]
	}
	//always not null
	role := *new([]string)
	roleS := ""
	if source.Role != nil && len(source.Role) > 0 {
		role = source.Role
		roleS = role[0]
	}

	intro := source.Intro
	if intro == "" {
		intro = aliasS + " " + roleS
	}

	return &model.Video{
		FindNo:       strings.ToUpper(strings.ReplaceAll(strings.ReplaceAll(source.Bangumi, "-", ""), "_", "")),
		Bangumi:      strings.ToUpper(source.Bangumi),
		Type:         source.Type,
		Systematics:  source.Systematics,
		Sharpness:    source.Sharpness,
		Producer:     source.Producer,
		Language:     source.Language,
		Caption:      source.Caption,
		Intro:        intro,
		Alias:        alias,
		Role:         role,
		Director:     source.Director,
		Series:       source.Series,
		Tags:         source.Tags,
		Date:         source.Date,
		SourceHash:   source.SourceHash,
		Season:       MustString(source.Season, "1"),
		Episode:      MustString(source.Episode, "1"),
		TotalEpisode: MustString(source.TotalEpisode, "1"),
		Format:       MustString(source.Format, "2D"),
		Publisher:    source.Publisher,
		Length:       source.Length,
		MagnetLinks:  source.MagnetLinks,
		Uncensored:   source.Uncensored,
	}
}

// defaultUnfinished ...
func defaultUnfinished(name string) *model.Unfinished {
	_, file := filepath.Split(name)

	uncat := &model.Unfinished{
		Model:       model.Model{},
		Checksum:    "",
		Type:        "other",
		Relate:      "",
		Name:        file,
		Hash:        "",
		Sharpness:   "",
		Caption:     "",
		Encrypt:     false,
		Key:         "",
		M3U8:        "media.m3u8",
		SegmentFile: "media-%05d.ts",
		Sync:        false,
		Object:      new(model.VideoObject),
	}
	log.With("file", name).Info("calculate checksum")
	uncat.Checksum = model.Checksum(name)
	return uncat
}

func checkFileNotExist(path string) bool {
	_, e := os.Stat(path)
	if e != nil {
		if os.IsNotExist(e) {
			return true
		}
		return false
	}
	return false
}

// Run ...
func (info *Information) Run(ctx context.Context) {
	log.Info("information running")
	var vs []*VideoSource
	isDefault := true
	switch info.InfoType {
	case InfoTypeBSON:
		isDefault = false
		b, e := ioutil.ReadFile(info.Path)
		if e != nil {
			panic(e)
		}
		fixed := fixBson(b)
		reader := bytes.NewBuffer(fixed)
		e = LoadFrom(&vs, reader)
		if e != nil {
			panic(e)
		}
	case InfoTypeJSON:
		isDefault = false
		b, e := ioutil.ReadFile(info.Path)
		if e != nil {
			panic(e)
		}
		reader := bytes.NewBuffer(b)
		e = LoadFrom(&vs, reader)
		if e != nil {
			panic(e)
		}
	default:
	}
	if !isDefault {
		if vs == nil {
			log.Info("nil video source")
			return
		}
		log.With("filter", info.ProcList).Info("filter list")
		vs = filterProcList(vs, info.ProcList)
		failedSkip := atomic.NewBool(false)
		maxLimit := len(vs)
		for i := info.Start; i < maxLimit; i++ {
			select {
			case <-ctx.Done():
				return
			default:
				source := vs[i]
				v := video(source)
				if !failedSkip.Load() {
					if source.Poster != "" {
						v.PosterHash = source.Poster
					} else {
						if source.PosterPath != "" {
							source.PosterPath = filepath.Join(info.ResourcePath, source.PosterPath)
							if checkFileNotExist(source.PosterPath) {
								log.With("index", i, "bangumi", source.Bangumi).Info("poster not found")
							} else {
								e := info.PushTo(StepperAPI, APICallback(source, func(api *API, api2 *httpapi.HttpApi, v interface{}) (e error) {
									source := v.(*VideoSource)
									file, e := os.Open(source.PosterPath)
									if e != nil {
										return e
									}
									resolved, e := api2.Unixfs().Add(ctx, files.NewReaderFile(file))
									if e != nil {
										return e
									}
									_, e = addPosterHash(info.Seeder, source, resolved.String())
									if e != nil {
										failedSkip.Store(true)
										return e
									}

									return nil
								}))
								if e != nil {
									log.Error(e)
									continue
								}

							}
						}
					}

					if source.Thumb != "" {
						source.Thumb = filepath.Join(info.ResourcePath, source.Thumb)
						if checkFileNotExist(source.Thumb) {
							log.With("index", i, "bangumi", source.Bangumi).Info("thumb not found")
						} else {
							e := info.PushTo(StepperAPI, APICallback(source, func(api *API, api2 *httpapi.HttpApi, v interface{}) (e error) {
								source := v.(*VideoSource)
								file, e := os.Open(source.PosterPath)
								if e != nil {
									return e
								}
								resolved, e := api2.Unixfs().Add(ctx, files.NewReaderFile(file))
								if e != nil {
									return e
								}
								_, e = addThumbHash(info.Seeder, source, resolved.String())
								if e != nil {
									failedSkip.Store(true)
									return e
								}
								return nil
							}))
							if e != nil {
								log.Error(e)
								continue
							}
						}
					}
				}
				e := info.PushTo(StepperDatabase, DatabaseCallback(v, func(database *Database, eng *xorm.Engine, v interface{}) (e error) {
					return model.AddOrUpdateVideo(eng.NewSession(), v.(*model.Video))
				}))
				if e != nil {
					log.With("bangumi", v.Bangumi).Error(e)
				}
			}
		}
		log.Info("info end")

	}
	return
}

func addThumbHash(seed Seeder, source *VideoSource, hash string) (unf *model.Unfinished, e error) {
	unfinThumb := defaultUnfinished(source.Thumb)
	unfinThumb.Type = model.TypeThumb
	unfinThumb.Relate = source.Bangumi
	if source.Thumb != "" {
		unfinThumb.Hash = hash
		e = seed.PushTo(StepperDatabase, DatabaseCallback(unfinThumb, func(database *Database, eng *xorm.Engine, v interface{}) (e error) {
			return model.AddOrUpdateUnfinished(eng.NewSession(), v.(*model.Unfinished))
		}))
		if e != nil {
			return nil, e
		}
		return unfinThumb, nil
	}

	return nil, xerrors.New("no thumb")
}

func addPosterHash(seed Seeder, source *VideoSource, hash string) (unf *model.Unfinished, e error) {
	unfinPoster := defaultUnfinished(source.PosterPath)
	unfinPoster.Type = model.TypePoster
	unfinPoster.Relate = source.Bangumi

	if source.PosterPath != "" {
		unfinPoster.Hash = hash
		e = seed.PushTo(StepperDatabase, DatabaseCallback(unfinPoster, func(database *Database, eng *xorm.Engine, v interface{}) (e error) {
			return model.AddOrUpdateUnfinished(eng.NewSession(), v.(*model.Unfinished))
		}))
		if e != nil {
			return nil, e
		}
		return unfinPoster, nil
	}
	return nil, xerrors.New("no poster")
}

func addVideo(db *Database, video *model.Video) {

}

func filterProcList(sources []*VideoSource, filterList []string) (vs []*VideoSource) {
	if filterList == nil || len(filterList) == 0 {
		return sources
	}
	for _, source := range sources {
		for _, v := range filterList {
			if source.Bangumi == v {
				vs = append(vs, source)
			}
		}
	}
	return
}

// VideoSource ...
type VideoSource struct {
	Bangumi      string    `json:"bangumi"`       //番号 no
	SourceHash   string    `json:"source_hash"`   //原片hash
	Type         string    `json:"type"`          //类型：film，FanDrama
	Format       string    `json:"format"`        //输出：3D，2D
	VR           string    `json:"vr"`            //VR格式：左右，上下，平面
	Thumb        string    `json:"thumb"`         //缩略图
	Intro        string    `json:"intro"`         //简介 title
	Alias        []string  `json:"alias"`         //别名，片名
	VideoEncode  string    `json:"video_encode"`  //视频编码
	AudioEncode  string    `json:"audio_encode"`  //音频编码
	Files        []string  `json:"files"`         //存放路径
	HashFiles    []string  `json:"hash_files"`    //已上传Hash
	CheckFiles   []string  `json:"check_files"`   //Unfinished checksum
	Slice        bool      `json:"sliceAdd"`      //是否HLS切片
	Encrypt      bool      `json:"encrypt"`       //加密
	Key          string    `json:"key"`           //秘钥
	M3U8         string    `json:"m3u8"`          //M3U8名
	SegmentFile  string    `json:"segment_file"`  //ts切片名
	PosterPath   string    `json:"poster_path"`   //海报路径
	Poster       string    `json:"poster"`        //海报HASH
	ExtendList   []*Extend `json:"extend_list"`   //扩展信息
	Role         []string  `json:"role"`          //角色列表 stars
	Director     string    `json:"director"`      //导演
	Systematics  string    `json:"systematics"`   //分级
	Season       string    `json:"season"`        //季
	Episode      string    `json:"episode"`       //集数
	TotalEpisode string    `json:"total_episode"` //总集数
	Sharpness    string    `json:"sharpness"`     //清晰度
	Publish      string    `json:"publish"`       //发行日
	Date         string    `json:"date"`          //发行日
	Length       string    `json:"length"`        //片长
	Producer     string    `json:"producer"`      //制片商
	Series       string    `json:"series"`        //系列
	Tags         []string  `json:"tags"`          //标签
	Publisher    string    `json:"publisher"`     //发行商
	Language     string    `json:"language"`      //语言
	Caption      string    `json:"caption"`       //字幕
	MagnetLinks  []string  `json:"magnet_links"`  //磁链
	Uncensored   bool      `json:"uncensored"`    //有码,无码
}

// InformationOption ...
func informationOption(info *Information) Options {
	return func(seed Seeder) {
		seed.SetThread(StepperInformation, info)
	}
}
