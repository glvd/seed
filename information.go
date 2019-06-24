package seed

import (
	"bytes"
	"context"
	"github.com/go-xorm/xorm"
	shell "github.com/godcong/go-ipfs-restapi"
	"github.com/yinhevr/seed/model"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
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

// InfoStatus ...
type InfoStatus string

// TransferFlagNone ...
const (
	InfoStatusNone   InfoStatus = "none"
	InfoStatusVerify InfoStatus = "verify"
	InfoStatusAdd    InfoStatus = "add"
	InfoStatusUpdate InfoStatus = "update"
	InfoStatusDelete InfoStatus = "delete"
)

// information ...
type information struct {
	workspace  string
	shell      *shell.Shell
	maindb     *xorm.Engine
	unfinished map[string]*model.Unfinished
	from       InfoFlag
	status     InfoStatus
	path       string
	videos     []*model.Video
}

// Information ...
func Information(path string, from InfoFlag, status InfoStatus) Options {
	info := &information{
		path:   path,
		from:   from,
		status: status,
	}
	return informationOption(info)
}

// BeforeRun ...
func (info *information) BeforeRun(seed *Seed) {
	info.workspace = seed.Workspace
	info.unfinished = seed.Unfinished
	info.shell = seed.Shell
	info.maindb = seed.maindb
}

// AfterRun ...
func (info *information) AfterRun(seed *Seed) {
	seed.Videos = info.videos
}

func fixBson(s []byte) []byte {
	reg := regexp.MustCompile(`("_id")[ ]*[:][ ]*(ObjectId\(")[\w]{24}("\))[ ]*(,)[ ]*`)
	return reg.ReplaceAll(s, []byte(" "))
}

func video(source *VideoSource) (video *model.Video) {
	video = new(model.Video)
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
	video.FindNo = strings.ReplaceAll(strings.ReplaceAll(source.Bangumi, "-", ""), "_", "")
	video.Bangumi = strings.ToUpper(source.Bangumi)
	video.Intro = intro
	video.Alias = alias
	video.Role = role
	video.Director = source.Director
	video.Series = source.Series
	video.Tags = source.Tags
	video.Date = source.Date
	video.SourceHash = source.SourceHash
	video.Season = MustString(source.Season, "1")
	video.Episode = MustString(source.Episode, "1")
	video.TotalEpisode = MustString(source.TotalEpisode, "1")
	video.Format = MustString(source.Format, "2D")
	video.Publisher = source.Publisher
	video.Length = source.Length
	return
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
		case InfoFlagSQLite:

		}
	}

	if vs == nil {
		log.Info("no video to process")
		return
	}

	for _, s := range vs {
		select {
		case <-ctx.Done():
			return
		default:
			v := video(s)
			s.Thumb = filepath.Join(info.workspace, s.Thumb)
			unfinThumb := DefaultUnfinished(s.Thumb)
			unfinThumb.Type = model.TypeThumb
			unfinThumb.Relate = s.Bangumi
			thumb, e := addThumbHash(info.shell, s)
			if e != nil {
				log.Error(e)
			}
			if thumb != "" {
				unfinThumb.Hash = thumb
				v.Thumb = thumb
				info.unfinished[v.Thumb] = unfinThumb
			}
			s.PosterPath = filepath.Join(info.workspace, s.PosterPath)
			unfinPoster := DefaultUnfinished(s.PosterPath)
			unfinPoster.Type = model.TypePoster
			unfinPoster.Relate = s.Bangumi
			poster, e := addPosterHash(info.shell, s)
			if e != nil {
				log.Error(e)
			}

			if poster != "" {
				v.PosterHash = poster
				unfinPoster.Hash = poster
				info.unfinished[v.PosterHash] = unfinPoster
			}

			if s.Poster != "" {
				v.PosterHash = s.Poster
			}

			info.videos = append(info.videos, v)

		}
	}

	return
}
