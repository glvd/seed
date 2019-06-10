package seed

import (
	"github.com/go-xorm/xorm"
	"github.com/yinhevr/seed/model"
	"golang.org/x/xerrors"
)

// TransferFlag ...
type TransferFlag string

// TransferFlagNone ...
const TransferFlagNone TransferFlag = "none"

// TransferFlagUpdate ...
const TransferFlagUpdate TransferFlag = "update"

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
	TransferStatusNone   TransferFlag = "none"
	TransferFlagVerify   TransferFlag = "verify"
	TransferStatusAdd    TransferFlag = "add"
	TransferStatusUpdate TransferFlag = "update"
	TransferStatusDelete TransferFlag = "delete"
)

// transfer ...
type transfer struct {
	from   TransferFlag
	to     TransferFlag
	status TransferStatus
}

// Transfer ...
func Transfer(from, to string) Options {
	return func(seed *Seed) {

	}
}

// Run ...
func (transfer *transfer) Run() (e error) {
	var videos = new([]*model.Video)

	if e = model.DB().Find(videos); e != nil {
		return
	}

	if e = WriteJSON("video.json", videos); e != nil {
		return
	}

	return nil
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
