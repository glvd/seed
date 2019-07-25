package seed

import (
	"context"
	"strconv"
	"sync"

	"github.com/yinhevr/seed/model"
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
	//var vs []*model.Video
	switch content {
	case UpdateContentAll:
		log.Info("update all")
		fallthrough
	case UpdateContentHash:
		log.With("bangumi", video.Bangumi).Info("update hash")
		unfins := new([]*model.Unfinished)
		i, e := model.DB().Where("relate like ?", video.Bangumi+"%").FindAndCount(unfins)
		if e != nil {
			return nil, e
		}
		if i <= 0 {
			return nil, nil
		}
		var unfin *model.Unfinished
		for j := i; j > 0; j-- {
			unfin = (*unfins)[j-1]
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

		for j := i; j > 0; j-- {
			unfin = (*unfins)[j-1]
			log.With("checksum", unfin.Checksum, "relate", unfin.Relate, "type", unfin.Type, "sharpness", unfin.Sharpness).Infof("unfinished")
			if idx := NumberIndex(unfin.Relate); idx != -1 {
				if vs[idx] == nil {
					vs[idx] = video.Clone()
					vs[idx].Episode = strconv.Itoa(idx + 1)
				}

				switch unfin.Type {
				case model.TypeSlice:
					//slice sharpness > source sharpness
					vs[idx].Sharpness = unfin.Sharpness
					vs[idx].M3U8Hash = unfin.Hash
				case model.TypeVideo:
					vs[idx].Sharpness = MustString(vs[idx].Sharpness, unfin.Sharpness)
					//vs[idx].Sharpness = unfin.Sharpness
					vs[idx].SourceHash = unfin.Hash
				}
				continue
			}
			if vs[0] == nil {
				vs[0] = video.Clone()
			}
			switch unfin.Type {
			case model.TypeSlice:
				//slice sharpness > source sharpness
				vs[0].Sharpness = unfin.Sharpness
				vs[0].M3U8Hash = unfin.Hash
			case model.TypeVideo:
				vs[0].Sharpness = MustString(vs[0].Sharpness, unfin.Sharpness)
				vs[0].SourceHash = unfin.Hash
			}
		}
		log.Infof("total(%d),value:%+v", len(vs), vs)
	case UpdateContentInfo:
		log.With("bangumi", video.Bangumi).Info("update info")
		unfins := new([]*model.Unfinished)

		i, e := model.DB().Where("relate like ?", video.Bangumi+"%").FindAndCount(unfins)
		if e != nil {
			return nil, e
		}
		var unfin *model.Unfinished
		for j := i; j > 0; j-- {
			unfin = (*unfins)[j-1]
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
		//vs = make([]*model.Video, i)

		for j := i; j > 0; j-- {
			unfin = (*unfins)[j-1]
			if idx := NumberIndex(unfin.Relate); idx != -1 {
				if strconv.Itoa(idx) == video.Episode {
					video.Sharpness = unfin.Sharpness
					switch unfin.Type {
					case model.TypeSlice:
						video.M3U8Hash = unfin.Hash
					case model.TypeVideo:
						video.SourceHash = unfin.Hash
					}
				}
			} else {
				video.Sharpness = unfin.Sharpness
				switch unfin.Type {
				case model.TypeSlice:
					video.M3U8Hash = unfin.Hash
				case model.TypeVideo:
					video.SourceHash = unfin.Hash
				}
			}
		}
		vs = []*model.Video{video}
		log.Infof("total(%d),value:%+v", len(vs), vs)
	}

	//for _, v := range vs {
	//	if v == nil {
	//		continue
	//	}
	//	result = append(result, v)
	//}

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
			i, e := model.DB().Count(&model.Video{})
			if e != nil {
				return
			}
			for j := int64(0); j < i; j += 50 {
				videos, e := model.AllVideos(nil, 50, int(j))
				if e != nil {
					return
				}
				for _, video := range *videos {
					vs, e := doContent(video, u.content)
					if e != nil {
						continue
					}
					//u.videos[video.Bangumi] = video
					for _, v := range vs {
						if v == nil {
							continue
						}
						vc <- v
					}
				}
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
				for _, v := range vs {
					if v == nil {
						continue
					}
					vc <- v
				}
			}
		case UpdateMethodVideo:
			for _, video := range u.videos {
				vs, e := doContent(video, u.content)
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
			e := model.AddOrUpdateVideo(v)
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
func (u *update) BeforeRun(seed *Seed) {
	u.videos = seed.Videos
	u.unfinished = seed.Unfinished
}

// AfterRun ...
func (u *update) AfterRun(seed *Seed) {

}
