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

// UpdateMethodVideo ...
const UpdateMethodVideo UpdateMethod = "video"

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

func doContent(video *model.Video, content UpdateContent) (vs []*model.Video, e error) {
	switch content {
	case UpdateContentAll:
		fallthrough
	case UpdateContentHash:
		unfins := new([]*model.Unfinished)
		i, e := model.DB().Where("relate like ?", video.Bangumi).FindAndCount(unfins)
		if e != nil {
			return nil, e
		}
		var unfin *model.Unfinished
		for ; i > 0; i-- {
			unfin = (*unfins)[i-1]
			switch unfin.Type {
			//case model.TypeSlice:
			//	video.M3U8Hash = unfin.Hash
			case model.TypePoster:
				video.PosterHash = unfin.Hash
			case model.TypeThumb:
				video.ThumbHash = unfin.Hash
				//case model.TypeVideo:
				//	video.SourceHash = unfin.Hash
			}
		}
		vs = make([]*model.Video, i)

		for ; i > 0; i-- {
			unfin = (*unfins)[i-1]
			if idx := noIndex(unfin.Relate); idx != -1 {
				vs[idx] = video.Clone()
				switch unfin.Type {
				case model.TypeSlice:
					vs[idx].M3U8Hash = unfin.Hash
				case model.TypeVideo:
					vs[idx].SourceHash = unfin.Hash
				}
				continue
			}
			vs[0] = video.Clone()
			switch unfin.Type {
			case model.TypeSlice:
				vs[0].M3U8Hash = unfin.Hash
			case model.TypeVideo:
				vs[0].SourceHash = unfin.Hash
			}
		}

	case UpdateContentInfo:
		//old, e := model.FindVideo(nil, video.Bangumi)
		//if e != nil {
		//	return e
		//}
		//video.ID = old.ID
		//video.SourceHash = old.SourceHash
		//video.M3U8Hash = old.M3U8Hash
		//video.PosterHash = old.PosterHash
		//video.ThumbHash = old.ThumbHash
		//video.Version = old.Version

	}
	return nil, nil
}

// Run ...
func (u *update) Run(context.Context) {
	log.Info("update running")
	var e error
	var updateVideos []*model.Video
	switch u.method {
	case UpdateMethodAll:
		videos, e := model.AllVideos(nil, 0)
		if e != nil {
			return
		}
		for _, video := range *videos {
			vs, e := doContent(video, u.content)
			if e != nil {
				continue
			}
			//u.videos[video.Bangumi] = video
			updateVideos = append(updateVideos, vs...)
		}
	case UpdateMethodUnfinished:
		for _, unfin := range u.unfinished {

			video, b := u.videos[unfin.Relate]
			if !b {
				relate := onlyNo(unfin.Relate)
				video, b := u.videos[relate]
				if b {
					video.Clone()
				}

				video, e = model.FindVideo(nil, unfin.Relate)
				if e != nil {
					log.With("id", unfin.ID).Error(e)
					continue
				}
			}
			vs, e := doContent(video, u.content)
			if e != nil {
				log.With("id", unfin.ID).Error(e)
				continue
			}
			updateVideos = append(updateVideos, vs...)
			//u.videos[video.Bangumi] = video
		}
	case UpdateMethodVideo:
		for _, video := range u.videos {
			vs, e := doContent(video, u.content)
			if e != nil {
				log.With("video", video.Bangumi).Error(e)
				continue
			}
			updateVideos = append(updateVideos, vs...)
		}
	}

	if updateVideos == nil {
		log.Error("nil videos")
		return
	}

	u.wg.Add(1)
	go func() {
		defer u.wg.Done()
		for _, video := range updateVideos {
			if video == nil {
				continue
			}
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
