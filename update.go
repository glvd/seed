package seed

import (
	"context"
	"github.com/yinhevr/seed/model"
	"sync"
)

// UpdateContent ...
type UpdateContent string

// UpdateMethod ...
type UpdateMethod string

// UpdateMethodBefore ...
const UpdateMethodBefore UpdateMethod = "before"

// UpdateMethodAll ...
const UpdateMethodAll UpdateMethod = "all"

// UpdateMethodUnfinished ...
const UpdateMethodUnfinished UpdateMethod = "unfinished"

// UpdateStatusNone ...
const (
	UpdateContentNone   UpdateContent = "none"
	UpdateContentVerify UpdateContent = "verify"
	UpdateContentAll    UpdateContent = "all"
	UpdateContentInfo   UpdateContent = "Info"
	UpdateContentHash   UpdateContent = "hash"
	UpdateContentDelete UpdateContent = "delete"
)

type update struct {
	wg         *sync.WaitGroup
	videos     map[string]*model.Video
	unfinished map[string]*model.Unfinished
	method     UpdateMethod
	content    UpdateContent
}

// Update ...
func Update(method UpdateMethod, content UpdateContent) Options {
	update := &update{
		method:  method,
		content: content,
		wg:      &sync.WaitGroup{},
	}
	return updateOption(update)
}

// updateOption ...
func updateOption(update *update) Options {
	return func(seed *Seed) {
		seed.thread[StepperUpdate] = update
	}
}

func doContent(video *model.Video, content UpdateContent) (e error) {
	if e != nil {
		return e
	}
	return nil
}

// Run ...
func (u *update) Run(context.Context) {
	var e error
	switch u.method {
	case UpdateMethodAll:
		videos, e := model.AllVideos(nil, 0)
		if e != nil {
			return
		}
		for _, video := range *videos {
			e := doContent(video, u.content)
			if e != nil {
				continue
			}
			u.videos[video.Bangumi] = video
		}
	case UpdateMethodUnfinished:

		for _, unfin := range u.unfinished {
			video, b := u.videos[unfin.Relate]
			if !b {
				video, e = model.FindVideo(nil, unfin.Relate)
				if e != nil {
					log.With("id", unfin.ID).Error(e)
					continue
				}
			}

			e := doContent(video, u.content)
			if e != nil {
				log.With("id", unfin.ID).Error(e)
				continue
			}

			u.videos[video.Bangumi] = video
		}
	case UpdateMethodBefore:

	}

	if u.videos == nil {
		log.Error("nil videos")
		return
	}

	u.wg.Add(1)
	go func() {
		defer u.wg.Done()
		for _, video := range u.videos {
			e := model.AddOrUpdateVideo(video)
			if e != nil {
				log.Error(e)
				continue
			}
		}
	}()
	u.wg.Wait()
}

// BeforeRun ...
func (u *update) BeforeRun(seed *Seed) {
	u.videos = seed.Videos
	u.unfinished = seed.Unfinished
}

// AfterRun ...
func (u *update) AfterRun(seed *Seed) {

}
