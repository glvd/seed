package seed

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"github.com/go-xorm/xorm"
	"github.com/yinhevr/seed/model"
	"golang.org/x/xerrors"
	"io"
	"io/ioutil"
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
	from   TransferFlag
	to     TransferFlag
	status TransferStatus
	path   string
	reader io.Reader
	video  []*model.Video
}

// BeforeRun ...
func (transfer *transfer) BeforeRun(seed *Seed) {
	b, e := ioutil.ReadFile(transfer.path)
	if e != nil {
		return
	}

	//e = ioutil.WriteFile("tmp.json", fixFile(b), os.ModePerm)
	//if e != nil {
	//	return
	//}
	//all, e := ioutil.ReadAll(transfer.reader)
	//if e != nil {
	//	return
	//}
	fixed := fixFile(b)
	transfer.reader = bytes.NewBuffer(fixed)
}

func fixFile(s []byte) []byte {
	reg := regexp.MustCompile(`("_id")[ ]*[:][ ]*(ObjectId\(")[\w]{24}("\))[ ]*(,)[ ]*`)
	return reg.ReplaceAll(s, []byte(" "))
}

// AfterRun ...
func (transfer *transfer) AfterRun(seed *Seed) {
	seed.Video = transfer.video
}

// TransferOption ...
func TransferOption(t *transfer) Options {
	return func(seed *Seed) {
		seed.thread[StepperTransfer] = t
	}
}

// Transfer ...
func Transfer(reader io.Reader, from, to TransferFlag, status TransferStatus) Options {
	t := &transfer{
		reader: reader,
		from:   from,
		to:     to,
		status: status,
	}
	return TransferOption(t)
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
			transfer.video = append(transfer.video, video(s))
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
	video.Thumb = source.Thumb
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
