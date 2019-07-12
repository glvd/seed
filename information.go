package seed

import (
	"bytes"
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	shell "github.com/godcong/go-ipfs-restapi"
	"github.com/yinhevr/seed/model"
	"go.uber.org/atomic"
)

// InfoFlag ...
type InfoFlag string

// InfoFlagNone ...
const InfoFlagNone InfoFlag = "none"

// InfoFlagInfo ...
const InfoFlagInfo InfoFlag = "information"

// InfoFlagUpdate ...
const InfoFlagUpdate InfoFlag = "update"

// InfoFlagMysql ...
const InfoFlagMysql InfoFlag = "mysql"

// InfoFlagJSON ...
const InfoFlagJSON InfoFlag = "json"

// InfoFlagBSON ...
const InfoFlagBSON InfoFlag = "bson"

// InfoFlagSQLite ...
const InfoFlagSQLite InfoFlag = "sqlite"

// information ...
type information struct {
	workspace  string
	shell      *shell.Shell
	unfinished map[string]*model.Unfinished
	from       InfoFlag
	path       string
	thread     int
	list       []string
	videos     map[string]*model.Video
	moves      map[string]string
	maxLimit   int
}

// Information ...
func Information(path string, from InfoFlag, list ...string) Options {
	info := &information{
		path: path,
		from: from,
		list: list,
	}
	return informationOption(info)
}

// BeforeRun ...
func (info *information) BeforeRun(seed *Seed) {
	info.workspace = seed.Workspace
	info.videos = seed.Videos
	info.unfinished = seed.Unfinished
	info.shell = seed.Shell
	info.moves = seed.Moves
	info.maxLimit = seed.MaxLimit
}

// AfterRun ...
func (info *information) AfterRun(seed *Seed) {
	seed.Videos = info.videos
}

func fixBson(s []byte) []byte {
	reg := regexp.MustCompile(`("_id")[ ]*[:][ ]*(ObjectId\(")[\w]{24}("\))[ ]*(,)[ ]*`)
	return reg.ReplaceAll(s, []byte(" "))
}

func videoChan(source *VideoSource, v chan<- *model.Video) {
	v <- video(source)
}

func video(source *VideoSource) (video *model.Video) {

	//always not null
	alias := []string{}
	aliasS := ""
	if source.Alias != nil && len(source.Alias) > 0 {
		alias = source.Alias
		aliasS = alias[0]
	}
	//always not null
	role := []string{}
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
		FindNo:       strings.ReplaceAll(strings.ReplaceAll(source.Bangumi, "-", ""), "_", ""),
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

// Run ...
func (info *information) Run(ctx context.Context) {
	log.Info("information running")
	var vs []*VideoSource
	select {
	case <-ctx.Done():
	default:
		switch info.from {
		case InfoFlagBSON:
			b, e := ioutil.ReadFile(info.path)
			if e != nil {
				return
			}
			fixed := fixBson(b)
			reader := bytes.NewBuffer(fixed)
			e = LoadFrom(&vs, reader)
			if e != nil {
				log.Error(e)
				return
			}
		case InfoFlagJSON:
			b, e := ioutil.ReadFile(info.path)
			if e != nil {
				return
			}
			reader := bytes.NewBuffer(b)
			e = LoadFrom(&vs, reader)
			if e != nil {
				log.Error(e)
				return
			}
		case InfoFlagMysql:
			fallthrough
		case InfoFlagSQLite:
			if info.list == nil {
				videos, e := model.AllVideos(nil, 0)
				if e != nil {
					log.Error(e)
					return
				}
				for _, video := range *videos {
					info.videos[video.Bangumi] = video
				}
				return
			}

			for _, name := range info.list {
				video, e := model.FindVideo(nil, name)
				if e != nil {
					log.Error(e)
					continue
				}
				info.videos[video.Bangumi] = video
			}
			//all work was done
			return
		}
	}

	if vs == nil {
		log.Info("no videos to process")
		return
	}
	vs = filterList(vs, info.list)
	max := len(vs)
	if max > info.maxLimit {
		max = info.maxLimit
	}
	log.With("size", len(vs), "max", max).Info("video source")
	skipIPFS := atomic.NewBool(false)
	v1 := make(chan *model.Video)
	skips := make(chan int)
	go func(v1 chan<- *model.Video, skp chan<- int) {
		runner := 0
		for i, s := range vs {
			log.With("index", i, "bangumi", s.Bangumi).Info("add info")
			if runner > max {
				v1 <- nil
			}
			v := video(s)
			if !skipIPFS.Load() {
				if s.PosterPath != "" {
					s.PosterPath = filepath.Join(info.workspace, s.PosterPath)
					if s.Poster != "" {
						v.PosterHash = s.Poster
					} else {
						poster, e := addPosterHash(info.shell, s)
						if os.IsNotExist(e) {
							continue
						}
						if e != nil {
							log.Error(e)
							skipIPFS.Store(true)
						} else {
							v.PosterHash = poster.Hash
							info.moves[poster.Hash] = s.PosterPath
						}
					}
				}

				if s.Thumb != "" {
					s.Thumb = filepath.Join(info.workspace, s.Thumb)
					thumb, e := addThumbHash(info.shell, s)
					if os.IsNotExist(e) {
						continue
					}
					if e != nil {
						log.Error(e)
						skipIPFS.Store(true)
					} else {
						v.ThumbHash = thumb.Hash
						info.moves[thumb.Hash] = s.Thumb
					}
				}
				runner++
			}
			v1 <- v
		}
	}(v1, skips)

	for ; max > 0; max-- {
		select {
		case v := <-v1:
			if v == nil {
				continue
			}
			info.videos[v.Bangumi] = v
			log.With("bangumi", v.Bangumi).Info("add video")
			e := model.AddOrUpdateVideo(v)
			if e != nil {
				log.Error(e)
				continue
			}
		}
	}

	return
}

func filterList(sources []*VideoSource, list []string) (vs []*VideoSource) {
	if list == nil || len(list) <= 0 {
		return sources
	}
	for _, source := range sources {
		for _, v := range list {
			if source.Bangumi == v {
				vs = append(vs, source)
			}
		}
	}
	return
}
