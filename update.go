package seed

import (
	"context"
	"strconv"
	"sync"

	"github.com/glvd/seed/model"
	"github.com/go-xorm/xorm"
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
	// UpdateContentNone ...
	UpdateContentNone UpdateContent = "none"
	// UpdateContentVerify ...
	UpdateContentVerify UpdateContent = "verify"
	// UpdateContentAll ...
	UpdateContentAll UpdateContent = "all"
	// UpdateContentInfo ...
	UpdateContentInfo UpdateContent = "Info"
	// UpdateContentHash ...
	UpdateContentHash UpdateContent = "hash"
	// UpdateContentDelete ...
	UpdateContentDelete UpdateContent = "delete"
)

type update struct {
	wg         *sync.WaitGroup
	videos     map[string]*model.Video
	unfinished map[string]*model.Unfinished
	method     UpdateMethod
	content    UpdateContent
}

// Push ...
func (u *update) Push(interface{}) error {
	return nil
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
	return func(seed Seeder) {
		seed.SetThread(StepperUpdate, update)
	}
}

func parseInfo(video *model.Video, unfin *model.Unfinished) {
	switch unfin.Type {
	case model.TypePoster:
		video.PosterHash = unfin.Hash
	case model.TypeThumb:
		video.ThumbHash = unfin.Hash
	case model.TypeSlice:
		video.Sharpness = unfin.Sharpness
		video.M3U8Hash = unfin.Hash
	case model.TypeVideo:
		video.Sharpness = MustString(video.Sharpness, unfin.Sharpness)
		video.SourceHash = unfin.Hash
	}
}

func doContent(engine *xorm.Engine, video *model.Video, content UpdateContent) (vs []*model.Video, e error) {
	//var vs []*model.Video
	switch content {
	case UpdateContentAll:
		log.Info("update all")
		fallthrough
	case UpdateContentHash:
		log.With("bangumi", video.Bangumi).Info("update hash")
		unfins := new([]*model.Unfinished)
		i, e := engine.Where("relate = ?", video.Bangumi).Or("relate like ?", video.Bangumi+"-%").FindAndCount(unfins)
		if e != nil {
			return nil, e
		}
		if i <= 0 {
			return nil, nil
		}
		var unfin *model.Unfinished
		for j := i; j > 0; j-- {
			unfin = (*unfins)[j-1]
			parseInfo(video, unfin)
		}

		vs = make([]*model.Video, i)
		total := 1
		for j := i; j > 0; j-- {
			unfin = (*unfins)[j-1]
			log.With("checksum", unfin.Checksum, "relate", unfin.Relate, "type", unfin.Type, "sharpness", unfin.Sharpness).Infof("unfinished")
			if idx := NumberIndex(unfin.Relate); idx != -1 {
				if vs[idx] == nil {
					vs[idx] = video.Clone()
					vs[idx].Episode = strconv.Itoa(idx + 1)
					if total < idx+1 {
						total = idx + 1
					}
				}

				parseInfo(vs[idx], unfin)
				continue
			}
			if vs[0] == nil {
				vs[0] = video.Clone()
			}
			parseInfo(vs[0], unfin)
		}

		for i := range vs {
			if vs[i] != nil {
				vs[i].TotalEpisode = strconv.Itoa(total)
			}
		}

		log.Infof("total(%d),value:%+v", len(vs), vs)
	case UpdateContentInfo:
		log.With("bangumi", video.Bangumi).Info("update info")
		unfins := new([]*model.Unfinished)

		i, e := engine.Where("relate like ?", video.Bangumi+"%").FindAndCount(unfins)
		if e != nil {
			return nil, e
		}
		var unfin *model.Unfinished
		for j := i; j > 0; j-- {
			unfin = (*unfins)[j-1]
			parseInfo(video, unfin)

		}

		for j := i; j > 0; j-- {
			unfin = (*unfins)[j-1]
			if idx := NumberIndex(unfin.Relate); idx != -1 {
				if strconv.Itoa(idx) == video.Episode {
					parseInfo(video, unfin)
				}
			} else {
				parseInfo(video, unfin)
			}
		}
		vs = []*model.Video{video}
		log.Infof("total(%d),value:%+v", len(vs), vs)
	}

	return vs, nil
}

// Run ...
func (u *update) Run(context.Context) {
	log.Info("update running")
	var e error
	videoChan := make(chan *model.Video, 10)
	go func(vc chan<- *model.Video) {
		switch u.method {
		case UpdateMethodAll:

			//i, e := model.DB().Count(&model.Video{})
			//if e != nil {
			//	return
			//}
			//for j := int64(0); j < i; j += 50 {
			//	videos, e := model.AllVideos(nil, 50, int(j))
			//	if e != nil {
			//		return
			//	}
			//	for _, video := range *videos {
			//		vs, e := doContent(video, u.content)
			//		if e != nil {
			//			continue
			//		}
			//		//u.videos[video.Bangumi] = video
			//		for _, v := range vs {
			//			if v == nil {
			//				continue
			//			}
			//			vc <- v
			//		}
			//	}
			//}
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
				vs, e := doContent(nil, video, u.content)
				if e != nil {
					log.With("id", unfin.ID).Error(e)
					continue
				}
				for _, v := range vs {
					if v == nil {
						continue
					}
					vc <- v
				}
			}
		case UpdateMethodVideo:
			for _, video := range u.videos {
				vs, e := doContent(nil, video, u.content)
				if e != nil {
					log.With("video", video.Bangumi).Error(e)
					continue
				}
				for _, v := range vs {
					if v == nil {
						continue
					}
					vc <- v
				}
			}
		}
		vc <- nil
	}(videoChan)

	for {
		select {
		case v := <-videoChan:
			if v == nil {
				goto END
			}
			log.With("bangumi", v.Bangumi, "m3u8_hash", v.M3U8Hash).Info("update")
			e := model.AddOrUpdateVideo(nil, v)
			if e != nil {
				log.Error(e)
				continue
			}
		}
	}

END:
	log.Info("update end")
}

// BeforeRun ...
func (u *update) BeforeRun(seed Seeder) {

}

// AfterRun ...
func (u *update) AfterRun(seed Seeder) {

}
