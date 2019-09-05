package task

import (
	"os"

	"github.com/glvd/seed"
	"github.com/glvd/seed/model"
	"github.com/go-xorm/xorm"
	cmd "github.com/godcong/go-ffmpeg-cmd"
	httpapi "github.com/ipfs/go-ipfs-http-client"
)

// VideoProcess ...
type VideoProcess struct {
	Path  string
	Scale int64
	Skip  []interface{}
}

// CallTask ...
func (v *VideoProcess) CallTask(seeder seed.Seeder, task *seed.Task) error {
	select {
	case <-seeder.Context().Done():
		return nil
	default:
		files := seed.GetFiles(v.Path)
		for _, f := range files {
			if seed.IsVideo(f) {
				call := videoCall{
					path:  f,
					scale: v.Scale,
					skip:  v.Skip,
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

// NewVideoProcess ...
func NewVideoProcess() *VideoProcess {
	path := os.TempDir()
	return &VideoProcess{
		Path: path,
	}
}

// Task ...
func (v *VideoProcess) Task() *seed.Task {
	return seed.NewTask(v)
}

type videoCall struct {
	path  string
	scale int64
	skip  []interface{}
}

// Call ...
func (call *videoCall) Call(process *seed.Process) (e error) {
	u := defaultUnfinished(call.path)
	u.Type = model.TypeVideo
	e = process.PushTo(seed.APICallback(u, func(api *seed.API, ipapi *httpapi.HttpApi, v interface{}) (e error) {
		u := v.(model.Unfinished)
		resolved, e := seed.AddFile(api, call.path)
		if e != nil {
			return e
		}
		u.Hash = model.PinHash(resolved)
		return api.PushTo(seed.DatabaseCallback(u, func(database *seed.Database, eng *xorm.Engine, v interface{}) (e error) {
			return model.AddOrUpdateUnfinished(eng.NewSession(), v.(*model.Unfinished))
		}))
	}))
	if e != nil {
		return e
	}

	e = process.PushTo(seed.SliceCall(call.path, u.Clone(), func(slice *seed.Slice, sa *cmd.SplitArgs, v interface{}) (e error) {
		//TODO:
		return nil
	}))
	if e != nil {
		return e
	}

}
