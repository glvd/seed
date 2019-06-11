package seed

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"github.com/go-xorm/xorm"
	shell "github.com/godcong/go-ipfs-restapi"
	"github.com/yinhevr/seed/model"
	"golang.org/x/xerrors"
	"io"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
)

// TransferFlag ...
type TransferFlag string

// TransferFlagNone ...
const TransferFlagNone TransferFlag = "none"

// TransferFlagUpdate ...
const TransferFlagUpdate TransferFlag = "updateAppHash"

// TransferFlagMysql ...
const TransferFlagMysql TransferFlag = "mysql"

// TransferFlagJSON ...
const TransferFlagJSON TransferFlag = "json"

// TransferFlagSQLite ...
const TransferFlagSQLite TransferFlag = "sqlite"

// TransferStatus ...
type TransferStatus string

// TransferFlagNone ...
const (
	TransferStatusNone   TransferStatus = "none"
	TransferFlagVerify   TransferStatus = "verify"
	TransferStatusAdd    TransferStatus = "add"
	TransferStatusUpdate TransferStatus = "updateAppHash"
	TransferStatusDelete TransferStatus = "delete"
)

// transfer ...
type transfer struct {
	shell      *shell.Shell
	unfinished map[string]*model.Unfinished
	from       TransferFlag
	to         TransferFlag
	status     TransferStatus
	path       string
	reader     io.Reader
	video      []*model.Video
}

// BeforeRun ...
func (transfer *transfer) BeforeRun(seed *Seed) {
	b, e := ioutil.ReadFile(transfer.path)
	if e != nil {
		return
	}
	transfer.shell = seed.Shell
	fixed := fixFile(b)
	transfer.reader = bytes.NewBuffer(fixed)

	transfer.unfinished = seed.Unfinished
	if transfer.unfinished == nil {
		transfer.unfinished = make(map[string]*model.Unfinished)
	}

}

func fixFile(s []byte) []byte {
	reg := regexp.MustCompile(`("_id")[ ]*[:][ ]*(ObjectId\(")[\w]{24}("\))[ ]*(,)[ ]*`)
	return reg.ReplaceAll(s, []byte(" "))
}

// AfterRun ...
func (transfer *transfer) AfterRun(seed *Seed) {
	seed.Video = transfer.video
	seed.Unfinished = transfer.unfinished
}

// TransferOption ...
func TransferOption(t *transfer) Options {
	return func(seed *Seed) {
		seed.thread[StepperTransfer] = t
	}
}

// Transfer ...
func Transfer(path string, from, to TransferFlag, status TransferStatus) Options {
	t := &transfer{
		path:   path,
		from:   from,
		to:     to,
		status: status,
	}
	return TransferOption(t)
}

func addThumbHash(tr *transfer, source *VideoSource) (string, error) {
	if source.Thumb != "" {
		abs, e := filepath.Abs(source.Thumb)
		if e != nil {
			return "", e
		}
		object, e := tr.shell.AddFile(abs)
		if e != nil {
			return "", e
		}
		return object.Hash, nil
	}
	return "", xerrors.New("no thumb")
}

func addPosterHash(tr *transfer, source *VideoSource) (string, error) {
	if source.PosterPath != "" {
		abs, e := filepath.Abs(source.PosterPath)
		if e != nil {
			return "", e
		}
		object, e := tr.shell.AddFile(abs)
		if e != nil {
			return "", e
		}
		return object.Hash, nil
	}
	if source.Poster != "" {
		return source.Poster, nil
	}
	return "", xerrors.New("no poster")
}

// Run ...
func (transfer *transfer) Run(ctx context.Context) {
	log.Info("transfer running")
	select {
	case <-ctx.Done():
	default:
		var vs []*VideoSource
		e := LoadFrom(&vs, transfer.reader)
		if e != nil {
			log.Error(e)
			return
		}
		for _, s := range vs {

			v := video(s)
			unfinThumb := DefaultUnfinished(s.Thumb)
			unfinThumb.Type = model.TypeThumb
			unfinThumb.Relate = s.Bangumi
			thumb, e := addThumbHash(transfer, s)
			if e != nil {
				log.Error(e)
			}
			v.Thumb = thumb

			unfinPoster := DefaultUnfinished(s.Poster)
			unfinPoster.Type = model.TypePoster
			unfinPoster.Relate = s.Bangumi
			poster, e := addPosterHash(transfer, s)
			if e != nil {
				log.Error(e)
			}
			v.PosterHash = poster
			transfer.video = append(transfer.video, v)
		}
	}

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
	//video.Type = source.Type
	//video.Format = source.Format
	//video.VR = source.VR
	//video.Thumb = source.Thumb
	video.Intro = intro
	video.Alias = alias
	video.Role = role
	video.Director = source.Director
	//video.Language = source.Language
	//video.Caption = source.Caption
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

// LoadFrom ...
func LoadFrom(vs *[]*VideoSource, reader io.Reader) (e error) {
	dec := json.NewDecoder(bufio.NewReader(reader))

	return dec.Decode(vs)
}

// TransferTo ...
func TransferTo(eng *xorm.Engine, limit int) (e error) {
	i, e := model.DB().Count(&model.Video{})
	if e != nil || i <= 0 {
		return e
	}
	if limit == 0 {
		limit = 10
	}
	for x := 0; x <= int(i); x += limit {
		var videos []*model.Video
		if e = model.DB().Limit(limit, x).Find(&videos); e != nil {
			return xerrors.Errorf("transfer error with:%d,%+v", x, e)
		}
		for _, v := range videos {
			log.Info("get:", v.Bangumi)
		}
		insert, e := eng.Insert(videos)
		log.Info(insert, e)

	}
	return nil
}
