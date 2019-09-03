package seed

import (
	"errors"
	"strconv"

	"github.com/glvd/seed/model"
	"github.com/go-xorm/xorm"
)

// UpdateContent ...
type UpdateContent string

// UpdateMethod ...
type UpdateMethod string

// UpdateMethod ...
const (
	UpdateMethodVideo UpdateMethod = "video"
	// UpdateMethodAll ...
	UpdateMethodAll UpdateMethod = "all"
	// UpdateMethodUnfinished ...
	UpdateMethodUnfinished UpdateMethod = "unfinished"
)

// UpdateStatus ...
const (
	// UpdateContentAll ...
	UpdateContentAll UpdateContent = "all"
	// UpdateContentInfo ...
	UpdateContentInfo UpdateContent = "info"
	// UpdateContentHash ...
	UpdateContentHash UpdateContent = "hash"
)

// Update ...
type Update struct {
	//*Thread
	cb         chan UpdateCaller
	updateFunc map[UpdateMethod]func(*Update)
	//method  UpdateMethod
	//content UpdateContent
}

// CallTask ...
func (u *Update) CallTask(seeder Seeder, task *Task) error {

	select {
	case <-seeder.Context().Done():
		return nil
		//TODO
	}

	return nil
}

// UpdateCallFunc ...
type UpdateCallFunc func(u *Update, f *xorm.Engine) error

type updateCall struct {
	cb       UpdateCallFunc
	database *xorm.Engine
	filter   []interface{}
}

// Call ...
func (uc *updateCall) Call(u *Update) error {
	return uc.cb(u, uc.database)
}

// UpdateCall ...
func UpdateCall(engine *xorm.Engine) (Stepper, UpdateCaller) {
	return StepperUpdate, &updateCall{
		cb:       callUpdate,
		database: engine,
	}
}

func callUpdate(u *Update, engine *xorm.Engine) error {
	return nil
}

// Push ...
func (u *Update) Push(v interface{}) error {
	return u.push(v)
}

// Task ...
func (u *Update) Task() *Task {
	return NewTask(u)
}

// NewUpdate ...
func NewUpdate() *Update {
	update := &Update{
		//method:  method,
		//content: content,
		//Thread: NewThread(),
	}
	return update
}

//updateOption ...
//func updateOption(update *Update) Options {
//	return func(seed Seeder) {
//		seed.SetBaseThread(StepperUpdate, update)
//	}
//}

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
		log.Info("Update all")
		fallthrough
	case UpdateContentHash:
		log.With("bangumi", video.Bangumi).Info("Update hash")
		unfins := new([]*model.Unfinished)
		i, e := engine.Where("relate = ?", video.Bangumi).Or("relate like ?", video.Bangumi+"@%").FindAndCount(unfins)
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
		log.With("bangumi", video.Bangumi).Info("Update info")
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
//func (u *Update) Run(ctx context.Context) {
//	log.Info("update running")
//UpdateEnd:
//	for {
//		select {
//		case <-ctx.Done():
//			break UpdateEnd
//		case cb := <-u.cb:
//			if cb == nil {
//				break UpdateEnd
//			}
//			u.SetState(StateRunning)
//			e := cb.Call(u)
//			if e != nil {
//				log.Error(e)
//			}
//		case <-time.After(30 * time.Second):
//			log.Info("update time out")
//			u.SetState(StateWaiting)
//		}
//	}
//	close(u.cb)
//	u.Finished()
//}

func (u *Update) push(v interface{}) error {
	if cb, b := v.(UpdateCaller); b {
		u.cb <- cb
		return nil
	}
	return errors.New("not update callback")
}
