package seed

import (
	"github.com/yinhevr/seed/model"
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
