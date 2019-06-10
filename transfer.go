package seed

import (
	"context"
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

// BeforeRun ...
func (transfer transfer) BeforeRun(seed *Seed) {

}

// AfterRun ...
func (transfer transfer) AfterRun(seed *Seed) {

}

// TransferOption ...
func TransferOption(t *transfer) Options {
	return func(seed *Seed) {
		seed.thread[StepperTransfer] = t
	}
}

// Transfer ...
func Transfer(from, to TransferFlag, status TransferStatus) Options {
	t := &transfer{
		from:   from,
		to:     to,
		status: status,
	}
	return TransferOption(t)
}

// Run ...
func (transfer *transfer) Run(ctx context.Context) {

}

// LoadFrom ...
func LoadFrom() {

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
