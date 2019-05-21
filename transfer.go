package seed

import (
	"github.com/go-xorm/xorm"
	"github.com/yinhevr/seed/model"
	"golang.org/x/xerrors"
)

// Transfer ...
func Transfer() (e error) {
	var videos = new([]*model.Video)

	if e = model.DB().Find(videos); e != nil {
		return
	}

	if e = WriteJSON("video.json", videos); e != nil {
		return
	}

	return nil
}

// TransferMysql ...
func TransferMysql(eng *xorm.Engine, limit int) (e error) {
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

// TransferOther ...
func TransferOther(name string, limit int) (e error) {
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
