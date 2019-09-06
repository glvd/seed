package task

import (
	"strconv"

	"github.com/glvd/seed"
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

// Update update info from the same db
type Update struct {
	Limit   int
	Include []interface{}
	Exclude []interface{}
}

// Task ...
func (u *Update) Task() *seed.Task {
	return seed.NewTask(u)
}

// NewUpdate ...
func NewUpdate() *Update {
	update := &Update{
		Limit: DefaultLimit,
	}
	return update
}

// CallTask ...
func (u *Update) CallTask(seeder seed.Seeder, task *seed.Task) error {
	select {
	case <-seeder.Context().Done():
		return nil
	default:
		return u.call(seeder)
	}
}

func (u *Update) call(seeder seed.Seeder) error {
	c := &updateCall{
		Limit:   u.Limit,
		Include: u.Include,
		Exclude: u.Exclude,
	}
	return seeder.PushTo(seed.StepperDatabase, c)
}

var _ seed.DatabaseCaller = &updateCall{}

type updateCall struct {
	Limit   int
	Include []interface{}
	Exclude []interface{}
}

// Call ...
func (u *updateCall) Call(database *seed.Database, eng *xorm.Engine) (e error) {
	session := eng.NewSession()
	if u.Include != nil {
		session = session.In("bangumi", u.Include...)
	}
	if u.Exclude != nil {
		session = session.NotIn("bangumi", u.Exclude)
	}
	v := make(chan *model.Video)

	go func(s *xorm.Session, video chan<- *model.Video) {
		defer s.Close()
		rows, e := session.Rows(new(model.Video))
		if e != nil {
			log.Error(e)
			return
		}
		count := 0
		defer func() {
			video <- nil
		}()
		for rows.Next() {
			v := new(model.Video)
			e := rows.Scan(v)
			if e != nil {
				log.Error(e)
				return
			}
			video <- v
			count++
			log.With("video", *v).Info("video info")
		}
		log.With("count", count).Info("done")
	}(session, v)

VideoEnd:
	for {
		select {
		case video := <-v:
			if video == nil {
				break VideoEnd
			}
			allUnfinished, e := model.AllUnfinished(eng.Where("relate = ?", video.Bangumi).Or("relate like ?", video.Bangumi+"@%"), 0)
			if e != nil {
				return e
			}
			videos := updateContentAll(video, *allUnfinished)
			for _, newVideo := range videos {
				e := model.AddOrUpdateVideo(eng.NewSession(), newVideo)
				if e != nil {
					log.Error(e)
				}
			}
		}
	}
	log.Info("update end")
	return nil
}

func parseTypeHash(u *model.Unfinished) func(video *model.Video) {
	switch u.Type {
	case model.TypeSlice:
		return func(video *model.Video) {
			video.Sharpness = u.Sharpness
			video.M3U8 = u.Hash
		}
	case model.TypeVideo:
		return func(video *model.Video) {
			video.Sharpness = seed.MustString(video.Sharpness, u.Sharpness)
			video.SourceHash = u.Hash
		}
	}
	return func(video *model.Video) {

	}
}

func parseTypeInfo(u *model.Unfinished) func(video *model.Video) {
	switch u.Type {
	case model.TypePoster:
		return func(video *model.Video) {
			video.PosterHash = u.Hash
		}
	case model.TypeThumb:
		return func(video *model.Video) {
			video.ThumbHash = u.Hash
		}
	}
	return func(video *model.Video) {

	}
}

func calcTotal(unfinishs []*model.Unfinished, fn func(u *model.Unfinished)) int {
	total := 1
	for _, u := range unfinishs {
		fn(u)
		if idx := seed.NumberIndex(u.Relate); idx != -1 {
			if total < idx+1 {
				total = idx + 1
			}
		}
	}
	return total
}

func updateContentAll(source *model.Video, unfinishs []*model.Unfinished) []*model.Video {
	size := len(unfinishs)
	//do nothing
	if unfinishs == nil || size == 0 {
		return []*model.Video{}
	}
	total := calcTotal(unfinishs, func(u *model.Unfinished) {
		parseTypeInfo(u)(source)
	})
	videos := make([]*model.Video, total)
	for _, u := range unfinishs {
		log.With("checksum", u.Checksum, "relate", u.Relate, "type", u.Type, "sharpness", u.Sharpness).Infof("unfinished")
		if idx := seed.NumberIndex(u.Relate); idx != -1 {
			if videos[idx] == nil {
				videos[idx] = source.Clone()
				videos[idx].Episode = strconv.Itoa(idx + 1)
				videos[idx].TotalEpisode = strconv.Itoa(total)
			}
			parseTypeHash(u)(videos[idx])
			continue
		}
		if videos[0] == nil {
			videos[0] = source.Clone()
		}
		parseTypeHash(u)(videos[0])
	}
	return videos
}
func updateContentInfo(source *model.Video, unfinishs []*model.Unfinished) []*model.Video {
	size := len(unfinishs)
	//do nothing
	if unfinishs == nil || size == 0 {
		return []*model.Video{}
	}
	_ = calcTotal(unfinishs, func(u *model.Unfinished) {
		parseTypeInfo(u)(source)
	})
	return []*model.Video{source}
}
