package task

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/glvd/seed"
	"github.com/glvd/seed/model"
	"github.com/glvd/seed/old"
	"github.com/go-xorm/xorm"
)

// TransferStatus ...
type TransferStatus string

// TransferFlagNone ...
const (
	// TransferStatusNone ...
	TransferStatusNone TransferStatus = "none"
	// TransferStatusVerify ...
	TransferStatusVerify TransferStatus = "verify"
	// TransferStatusToJSON ...
	TransferStatusToJSON TransferStatus = "json"
	// TransferStatusFromOther ...
	TransferStatusFromOther TransferStatus = "other"
	// TransferStatusFromOld ...
	TransferStatusFromOld TransferStatus = "old"
	// TransferStatusUpdate ...
	TransferStatusUpdate TransferStatus = "Update"
	// TransferStatusDelete ...
	TransferStatusDelete TransferStatus = "delete"
)

// TransferFlag ...
type TransferFlag string

// TransferFlagSQL ...
const TransferFlagSQL TransferFlag = "sql"

// TransferFlagJSON ...
const TransferFlagJSON TransferFlag = "json"

// Transfer transfer info from other json,database
type Transfer struct {
	flag     TransferFlag
	database *xorm.Engine
	path     string
	Status   TransferStatus
	Limit    int
	Before   *time.Time
	After    *time.Time
}

// Task ...
func (t *Transfer) Task() *seed.Task {
	return seed.NewTask(t)
}

// CallTask ...
func (t *Transfer) CallTask(seeder seed.Seeder, task *seed.Task) error {
	select {
	case <-seeder.Context().Done():
		return nil
	default:
		switch t.flag {
		case TransferFlagJSON:
			t := &jsonTransfer{
				flag:   t.flag,
				status: t.Status,
				path:   t.path,
			}
			e := seeder.PushTo(seed.StepperDatabase, t)
			if e != nil {
				return e
			}
		case TransferFlagSQL:
			t := &dbTransfer{
				database: t.database,
				status:   t.Status,
			}
			e := seeder.PushTo(seed.StepperDatabase, t)
			if e != nil {
				return e
			}
		}
	}
	return nil
}

// NewJSONTransfer ...
func NewJSONTransfer(path string) *Transfer {
	t := &Transfer{
		flag:  TransferFlagJSON,
		path:  path,
		Limit: DefaultLimit,
	}
	return t
}

type jsonTransfer struct {
	flag   TransferFlag
	status TransferStatus
	path   string
}

// Call ...
func (j *jsonTransfer) Call(database *seed.Database, eng *xorm.Engine) (e error) {
	switch j.status {
	case TransferStatusFromOther:
		fallthrough
	case TransferStatusFromOld:
		return errors.New("inputted a wrong status type")
	case TransferStatusToJSON:
		file, e := os.OpenFile(j.path, os.O_RDWR|os.O_TRUNC|os.O_SYNC|os.O_CREATE, os.ModePerm)
		if e != nil {
			return e
		}
		enc := json.NewEncoder(file)
		videos, e := model.AllVideos(eng.NoCache(), 0)
		if e != nil {
			return e
		}
		e = enc.Encode(*videos)
		if e != nil {
			return e
		}
	}
	return nil
}

// NewDBTransfer ...
func NewDBTransfer(db *xorm.Engine) *Transfer {
	t := &Transfer{
		flag:     TransferFlagSQL,
		database: db,
		Status:   TransferStatusFromOther,
		Limit:    DefaultLimit,
	}
	return t
}

type dbTransfer struct {
	database *xorm.Engine
	status   TransferStatus
	limit    int
}

// Call ...
func (d *dbTransfer) Call(database *seed.Database, eng *xorm.Engine) (e error) {
	switch d.status {
	case TransferStatusFromOther:
		e = copyUnfinished(eng, d.database, d.limit)
		if e != nil {
			return e
		}
		e = copyVideo(eng, d.database)
		if e != nil {
			return e
		}
	case TransferStatusFromOld:
		//now has no old data
	}
	return nil
}

func copyVideo(to *xorm.Engine, from *xorm.Engine) error {
	v := new(model.Video)

	rows, e := from.Rows(&model.Video{})
	if e != nil {
		return e
	}
	for rows.Next() {
		e := rows.Scan(v)
		if e != nil {
			log.Error(e)
		}

		e = model.AddOrUpdateVideo(to.NoCache(), v)
		if e != nil {
			log.With("bangumi", v.Bangumi).Error(e)
			continue
		}
	}

	return e
}

func insertOldToUnfinished(eng *xorm.Engine, ban string, obj *old.Object) error {
	hash := ""
	if obj.Link != nil {
		hash = obj.Link.Hash
	}
	unf := &model.Unfinished{
		Checksum:    hash,
		Type:        model.TypeSlice,
		Relate:      ban,
		Name:        ban,
		Hash:        hash,
		Sharpness:   "",
		Caption:     "",
		Encrypt:     false,
		Key:         "",
		M3U8:        "",
		SegmentFile: "",
		Sync:        false,
		Object:      ObjectFromOld(obj),
	}
	return model.AddOrUpdateUnfinished(eng.NoCache(), unf)

}

// ObjectFromOld ...
func ObjectFromOld(obj *old.Object) *model.VideoObject {
	return &model.VideoObject{
		Links: obj.Links,
		Link:  obj.Link,
	}
}

func oldToVideo(source *old.Video) *model.Video {
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
		FindNo:       strings.ToUpper(strings.ReplaceAll(strings.ReplaceAll(source.Bangumi, "-", ""), "_", "")),
		Bangumi:      strings.ToUpper(source.Bangumi),
		Intro:        intro,
		Alias:        alias,
		Role:         role,
		Season:       "1",
		Episode:      "1",
		TotalEpisode: "1",
		Format:       "2D",
	}

}

func transferFromOld(engine *xorm.Engine) (e error) {
	videos := old.AllVideos(engine)
	log.With("size", len(videos)).Info("videos")
	for _, v := range videos {
		obj := old.GetObject(v)
		e := insertOldToUnfinished(engine, v.Bangumi, obj)
		if e != nil {
			log.With("bangumi", v.Bangumi).Error(e)
			continue
		}
		vd, e := model.FindVideo(nil, v.Bangumi)
		if e != nil {
			log.With("bangumi", v.Bangumi).Error(e)
			continue
		}

		log.With("bangumi", v.Bangumi, "v", vd).Info("v Update")
		if vd.ID == "" {
			vd = oldToVideo(v)
		}

		if strings.TrimSpace(vd.M3U8Hash) == "" && obj.Link != nil {
			log.With("hash:", obj.Link.Hash, "bangumi", v.Bangumi).Info("info")
			vd.M3U8Hash = obj.Link.Hash
			e = model.AddOrUpdateVideo(nil, vd)
			if e != nil {
				log.With("bangumi", v.Bangumi).Error(e)
				continue
			}
		} else {

		}

	}
	return e
}

func copyUnfinished(to *xorm.Engine, from *xorm.Engine, limit int) (e error) {
	i, e := from.Count(&model.Unfinished{})
	if e != nil {
		log.Error(e)
		return
	}
	if limit < 500 {
		limit = 500
	}
	u := new(model.Unfinished)
	for x := int64(0); x < i; x += int64(limit) {
		rows, e := from.Limit(limit, int(x)).Rows(&model.Unfinished{})
		if e != nil {
			return e
		}
		for rows.Next() {
			e := rows.Scan(u)
			if e != nil {
				log.Error(e)
			}
			u.ID = ""
			u.Version = 0
			e = model.AddOrUpdateUnfinished(to.NoCache(), u)
			log.With("checksum", u.Checksum, "type", u.Type, "relate", u.Relate, "error", e).Info("copy")
			if e != nil {
				log.Error(e)
				continue
			}
		}
	}

	log.Infof("unfinished(%d) done", i)
	return nil
}

func transferUpdate(engine *xorm.Engine) (e error) {
	fromList := new([]*model.Unfinished)
	e = engine.Find(fromList)
	if e != nil {
		return
	}
	for _, from := range *fromList {
		video, e := model.FindVideo(engine.Where("episode = ?", seed.NumberIndex(from.Relate)), seed.OnlyName(from.Relate))
		if e != nil {
			log.Error(e)
			continue
		}

		if from.Type == model.TypeSlice {
			video.Sharpness = seed.MustString(from.Sharpness, video.Sharpness)
			video.M3U8Hash = seed.MustString(from.Hash, video.M3U8Hash)
		} else if from.Type == model.TypeVideo {
			video.Sharpness = seed.MustString(from.Sharpness, video.Sharpness)
			video.SourceHash = seed.MustString(from.Hash, video.SourceHash)
		} else {

		}
		e = model.AddOrUpdateVideo(nil, video)
		if e != nil {
			log.With("bangumi", video.Bangumi, "index", video.Episode).Error(e)
			continue
		}
	}
	return e
}

func transferToJSON(engine *xorm.Engine, to string) (e error) {
	videos, e := model.AllVideos(engine.Where("m3u8_hash <> ?", ""), 0)
	if e != nil {
		return e
	}
	bytes, e := json.Marshal(videos)
	if e != nil {
		return e
	}
	file, e := os.OpenFile(to, os.O_CREATE|os.O_SYNC|os.O_RDWR, os.ModePerm)
	if e != nil {
		return e
	}
	defer file.Close()
	n, e := file.Write(bytes)
	log.With("video", len(*videos)).Infof("write(%d)", n)
	return e
}

// Run ...
func (t *Transfer) Run(ctx context.Context) {
	if t.flag == TransferFlagSQL {
		fromDB, e := model.InitSQLite3(t.path)
		if e != nil {
			log.Error(e)
			return
		}
		e = fromDB.Sync2(model.Video{})
		if e != nil {
			log.Error(e)
			return
		}
		switch t.Status {
		case TransferStatusFromOld:
			if err := transferFromOld(fromDB); err != nil {
				log.Error(err)
				return
			}
		//Update flag video flag other sqlite3
		case TransferStatusFromOther:

			//if err := transferFromOther(fromDB); err != nil {
			//	return
			//}
		//Update flag unfinished flag other sqlite3
		case TransferStatusUpdate:
			if err := transferUpdate(fromDB); err != nil {
				return
			}
		}
	} else if t.flag == TransferFlagJSON {
		switch t.Status {
		case TransferStatusToJSON:
			if err := transferToJSON(nil, t.path); err != nil {
				return
			}
		}
	}

}

// LoadFrom ...
func LoadFrom(vs *[]*VideoSource, reader io.Reader) (e error) {
	dec := json.NewDecoder(bufio.NewReader(reader))
	return dec.Decode(vs)
}

// TransferTo ...
func TransferTo(eng *xorm.Engine, limit int) (e error) {
	i, e := eng.Count(&model.Video{})
	if e != nil || i <= 0 {
		return e
	}
	if limit == 0 {
		limit = 10
	}
	for x := 0; x <= int(i); x += limit {
		var videos []*model.Video
		if e = eng.Limit(limit, x).Find(&videos); e != nil {
			return fmt.Errorf("transfer error with:%d,%+v", x, e)
		}
		for _, v := range videos {
			log.Info("get:", v.Bangumi)
		}
		insert, e := eng.Insert(videos)
		log.Info(insert, e)
	}

	return nil
}
