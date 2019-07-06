package seed

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"path/filepath"
	"strings"

	"github.com/go-xorm/xorm"
	shell "github.com/godcong/go-ipfs-restapi"
	"github.com/yinhevr/seed/model"
	"github.com/yinhevr/seed/old"
	"golang.org/x/xerrors"
)

// TransferStatus ...
type TransferStatus string

// TransferFlagNone ...
const (
	TransferStatusNone   TransferStatus = "none"
	TransferFlagVerify   TransferStatus = "verify"
	TransferStatusAdd    TransferStatus = "add"
	TransferStatusOld    TransferStatus = "old"
	TransferStatusUpdate TransferStatus = "update"
	TransferStatusDelete TransferStatus = "delete"
)

// transfer ...
type transfer struct {
	shell      *shell.Shell
	unfinished map[string]*model.Unfinished
	videos     map[string]*model.Video
	workspace  string
	from       InfoFlag
	to         InfoFlag
	status     TransferStatus
	path       string
	reader     io.Reader
}

// BeforeRun ...
func (transfer *transfer) BeforeRun(seed *Seed) {
	transfer.shell = seed.Shell
	transfer.workspace = seed.Workspace
	transfer.unfinished = seed.Unfinished
	transfer.videos = seed.Videos

}

// AfterRun ...
func (transfer *transfer) AfterRun(seed *Seed) {
	seed.Videos = transfer.videos
	seed.Unfinished = transfer.unfinished
}

// TransferOption ...
func TransferOption(t *transfer) Options {
	return func(seed *Seed) {
		seed.thread[StepperTransfer] = t
	}
}

// Transfer ...
func Transfer(path string, from InfoFlag, status TransferStatus) Options {
	t := &transfer{
		path:   path,
		from:   from,
		status: status,
	}
	return TransferOption(t)
}

func addThumbHash(shell *shell.Shell, source *VideoSource) (*model.Unfinished, error) {
	unfinThumb := defaultUnfinished(source.Thumb)
	unfinThumb.Type = model.TypeThumb
	unfinThumb.Relate = source.Bangumi
	if source.Thumb != "" {
		abs, e := filepath.Abs(source.Thumb)
		if e != nil {
			return nil, e
		}

		object, e := shell.AddFile(abs)
		if e != nil {
			return nil, e
		}

		unfinThumb.Hash = object.Hash
		e = model.AddOrUpdateUnfinished(unfinThumb)
		if e != nil {
			return nil, e
		}
		return unfinThumb, nil
	}

	return nil, xerrors.New("no thumb")
}

func addPosterHash(shell *shell.Shell, source *VideoSource) (*model.Unfinished, error) {
	unfinPoster := defaultUnfinished(source.PosterPath)
	unfinPoster.Type = model.TypePoster
	unfinPoster.Relate = source.Bangumi

	if source.PosterPath != "" {
		abs, e := filepath.Abs(source.PosterPath)
		if e != nil {
			return nil, e
		}
		object, e := shell.AddFile(abs)
		if e != nil {
			return nil, e
		}
		unfinPoster.Hash = object.Hash
		e = model.AddOrUpdateUnfinished(unfinPoster)
		if e != nil {
			return nil, e
		}
		return unfinPoster, nil
	}
	return nil, xerrors.New("no poster")
}

func insertOldToUnfinished(ban string, obj *old.Object) error {
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
	return model.AddOrUpdateUnfinished(unf)

}

// ObjectFromOld ...
func ObjectFromOld(obj *old.Object) *model.VideoObject {
	return &model.VideoObject{
		Links: obj.Links,
		Link:  obj.Link,
	}
}

// Run ...
func (transfer *transfer) Run(ctx context.Context) {
	switch transfer.from {
	case InfoFlagSQLite:
		if transfer.status == TransferStatusOld {
			objects := old.LoadFrom(transfer.path)
			log.With("size", len(objects)).Info("objects")
			for ban, obj := range objects {
				e := insertOldToUnfinished(ban, obj)
				if e != nil {
					log.With("bangumi", ban).Error(e)
					continue
				}
				vd, e := model.FindVideo(nil, ban)
				if e != nil || vd.ID == "" {
					log.With("bangumi", ban).Error(e)
					continue
				}
				log.With("bangumi", ban, "video", vd).Info("video update")
				if strings.TrimSpace(vd.M3U8Hash) == "" && obj.Link != nil {
					log.With("hash:", obj.Link.Hash, "bangumi", ban).Info("info")
					vd.M3U8Hash = obj.Link.Hash
					e = model.AddOrUpdateVideo(vd)
					if e != nil {
						log.With("bangumi", ban).Error(e)
						continue
					}
				} else {

				}

			}
		} else if transfer.status == TransferStatusUpdate {

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
