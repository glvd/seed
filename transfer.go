package seed

import (
	"github.com/go-xorm/xorm"
	"github.com/yinhevr/seed/model"
	"log"
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
func TransferMysql(eng *xorm.Engine) (e error) {

	i, e := model.DB().Count(&model.Video{})
	if e != nil {
		return e
	}
	video := model.Video{}
	for x := int64(0); x < i; x++ {
		b, e := model.DB().Limit(1, int(x)).Get(&video)
		log.Printf("%+v", video)
		if e != nil {
			return e
		}
		if !b {
			continue
		}
	}
	return nil
}
