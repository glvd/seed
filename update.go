package seed

import (
	"context"
	"github.com/yinhevr/seed/model"
	"sync"
)

type update struct {
	wg         *sync.WaitGroup
	video      []*model.Video
	unfinished map[string]*model.Unfinished
}

// Update ...
func Update() Options {
	update := &update{
		wg: &sync.WaitGroup{},
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
	if u.video == nil {
		log.Error("nil video")
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
		for _, video := range u.video {
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
	u.video = seed.Videos
	u.unfinished = seed.Unfinished
}

// AfterRun ...
func (u *update) AfterRun(seed *Seed) {

}
