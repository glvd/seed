package task

import (
	"os"

	"github.com/glvd/seed"
	"github.com/glvd/seed/model"
	"github.com/go-xorm/xorm"
	cmd "github.com/godcong/go-ffmpeg-cmd"
	httpapi "github.com/ipfs/go-ipfs-http-client"
)

// VideoSlice ...
type VideoSlice struct {
	Path     string
	SkipType []interface{}
	Filter   []string
}

// CallTask ...
func (v *VideoSlice) CallTask(seeder seed.Seeder, task *seed.Task) error {
	select {
	case <-seeder.Context().Done():
		return nil
	default:
		files := seed.GetFiles(v.Path)
		for _, f := range files {
			if seed.IsVideo(f) {
				call := &videoCall{
					path:     f,
					skipType: v.SkipType,
				}
				e := seeder.PushTo(seed.StepperProcess, call)
				if e != nil {
					log.Error(e)
					continue
				}
			}
		}
	}

	return nil
}

// NewVideoSlice ...
func NewVideoSlice() *VideoSlice {
	path := os.TempDir()
	return &VideoSlice{
		Path: path,
	}
}

// Task ...
func (v *VideoSlice) Task() *seed.Task {
	return seed.NewTask(v)
}

type videoCall struct {
	path     string
	skipType []interface{}
}

// Call ...
func (call *videoCall) Call(process *seed.Process) (e error) {
	u := defaultUnfinished(call.path)
	u.Type = model.TypeVideo
	f, err := cmd.FFProbeStreamFormat(call.path)
	if err != nil {
		return err
	}
	u.Sharpness = f.Resolution() + "P"
	u.Relate = seed.OnlyName(call.path)
	if !seed.SkipTypeVerify(u.Type, call.skipType...) {
		e = process.PushTo(seed.APICallback(u.Clone(), func(api *seed.API, ipapi *httpapi.HttpApi, v interface{}) (e error) {
			u := v.(*model.Unfinished)
			resolved, e := seed.AddFile(api, call.path)
			if e != nil {
				return e
			}
			u.Hash = model.PinHash(resolved)
			log.With("hash", u.Hash, "sharpness", u.Sharpness).Info("video")
			return api.PushTo(seed.DatabaseCallback(u, func(database *seed.Database, eng *xorm.Engine, v interface{}) (e error) {
				return model.AddOrUpdateUnfinished(eng.NoCache(), v.(*model.Unfinished))
			}))
		}))
		if e != nil {
			return e
		}
	}

	u.Type = model.TypeSlice
	e = process.PushTo(seed.SliceCall(call.path, u.Clone(), func(slice *seed.Slice, sa *cmd.SplitArgs, v interface{}) (e error) {
		u := v.(*model.Unfinished)
		return slice.PushTo(seed.APICallback(u.Clone(), func(api *seed.API, ipapi *httpapi.HttpApi, v interface{}) (e error) {
			u := v.(*model.Unfinished)
			resolved, e := seed.AddDir(api, sa.Output)
			if e != nil {
				return e
			}
			u.Hash = model.PinHash(resolved)
			log.With("hash", u.Hash, "sharpness", u.Sharpness).Info("slice")
			return api.PushTo(seed.DatabaseCallback(u, func(database *seed.Database, eng *xorm.Engine, v interface{}) (e error) {
				return model.AddOrUpdateUnfinished(eng.NoCache(), v.(*model.Unfinished))
			}))
		}))
	}))
	if e != nil {
		return e
	}
	return nil
}
