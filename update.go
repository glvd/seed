package seed

import (
	"context"
	"github.com/yinhevr/seed/model"
)

type update struct {
	video      []*model.Video
	unfinished map[string]*model.Unfinished
}

// Update ...
func Update() Options {
	update := &update{}
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

	for _, unfin := range u.unfinished {
		if err := model.AddOrUpdateUnfinished(unfin); err != nil {
			log.Error(err)
			continue
		}
	}

	for _, video := range u.video {
		e := model.AddOrUpdateVideo(video)
		if e != nil {
			log.Error(e)
			continue
		}
	}

}

// BeforeRun ...
func (u *update) BeforeRun(seed *Seed) {
	u.video = seed.Video
	u.unfinished = seed.Unfinished
}

// AfterRun ...
func (u *update) AfterRun(seed *Seed) {

}
