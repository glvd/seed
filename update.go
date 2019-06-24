package seed

import (
	"context"
	"github.com/yinhevr/seed/model"
	"sync"
)

// UpdateStatus ...
type UpdateStatus string

// UpdateStatusNone ...
const (
	UpdateStatusNone   UpdateStatus = "none"
	UpdateStatusVerify UpdateStatus = "verify"
	UpdateStatusAdd    UpdateStatus = "add"
	UpdateStatusUpdate UpdateStatus = "update"
	UpdateStatusDelete UpdateStatus = "delete"
)

type update struct {
	wg         *sync.WaitGroup
	videos     map[string]*model.Video
	unfinished map[string]*model.Unfinished
	status     UpdateStatus
}

// Update ...
func Update(status UpdateStatus) Options {
	update := &update{
		status: status,
		wg:     &sync.WaitGroup{},
	}
	return UpdateOption(update)
}

// UpdateOption ...
func UpdateOption(update *update) Options {
	return func(seed *Seed) {
		seed.thread[StepperUpdate] = update
	}
}

// Run ...
func (u *update) Run(context.Context) {
	if u.videos == nil {
		log.Error("nil videos")
		return
	}
	u.wg.Add(1)
	go func() {
		defer u.wg.Done()
		for _, unfin := range u.unfinished {
			if err := model.AddOrUpdateUnfinished(unfin); err != nil {
				log.Error(err)
				continue
			}
		}
	}()

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
