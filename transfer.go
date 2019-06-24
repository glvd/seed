package seed

import (
	"bufio"
	"context"
	"encoding/json"
	"github.com/go-xorm/xorm"
	shell "github.com/godcong/go-ipfs-restapi"
	"github.com/yinhevr/seed/model"
	"golang.org/x/xerrors"
	"io"
	"path/filepath"
)

// TransferStatus ...
type TransferStatus string

// TransferFlagNone ...
const (
	TransferStatusNone   TransferStatus = "none"
	TransferFlagVerify   TransferStatus = "verify"
	TransferStatusAdd    TransferStatus = "add"
	TransferStatusUpdate TransferStatus = "update"
	TransferStatusDelete TransferStatus = "delete"
)

// transfer ...
type transfer struct {
	shell      *shell.Shell
	unfinished map[string]*model.Unfinished
	workspace  string
	from       InfoFlag
	to         InfoFlag
	status     TransferStatus
	path       string
	reader     io.Reader
	video      []*model.Video
}

// BeforeRun ...
func (transfer *transfer) BeforeRun(seed *Seed) {
	transfer.shell = seed.Shell
	transfer.workspace = seed.Workspace
	transfer.unfinished = seed.Unfinished
	if transfer.unfinished == nil {
		transfer.unfinished = make(map[string]*model.Unfinished)
	}

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
func Transfer(path string, from, to InfoFlag, status TransferStatus) Options {
	t := &transfer{
		path:   path,
		from:   from,
		to:     to,
		status: status,
	}
	return TransferOption(t)
}

func addThumbHash(shell *shell.Shell, source *VideoSource) (string, error) {
	if source.Thumb != "" {
		abs, e := filepath.Abs(source.Thumb)
		if e != nil {
			return "", e
		}
		object, e := shell.AddFile(abs)
		if e != nil {
			return "", e
		}
		return object.Hash, nil
	}
	return "", xerrors.New("no thumb")
}

func addPosterHash(shell *shell.Shell, source *VideoSource) (string, error) {
	if source.PosterPath != "" {
		abs, e := filepath.Abs(source.PosterPath)
		if e != nil {
			return "", e
		}
		object, e := shell.AddFile(abs)
		if e != nil {
			return "", e
		}
		return object.Hash, nil
	}
	return "", xerrors.New("no poster")
}

// Run ...
func (transfer *transfer) Run(ctx context.Context) {
	log.Info("transfer running")

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
