package seed

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-xorm/xorm"
	files "github.com/ipfs/go-ipfs-files"
	httpapi "github.com/ipfs/go-ipfs-http-client"
	"golang.org/x/xerrors"

	"github.com/glvd/seed/model"
	cmd "github.com/godcong/go-ffmpeg-cmd"
)

// SliceCaller ...
type SliceCaller interface {
	Path() string
	Call(*Slice, string) error
}

// Slice ...
type Slice struct {
	Seeder
	Scale       int64
	SliceOutput string
	SkipType    []interface{}
	SkipExist   bool
	SkipSlice   bool
	cb          chan SliceCaller
}

// Push ...
func (s *Slice) Push(cb interface{}) error {
	if v, b := cb.(SliceCaller); b {
		s.cb <- v
		return nil
	}
	return xerrors.New("not slice callback")
}

// NewSlice ...
func NewSlice() *Slice {
	return &Slice{}
}

// Run ...
func (s *Slice) Run(context.Context) {
	log.Info("slice running")

	for {
		select {
		case v := <-s.cb:
			files := GetFiles(v.Path())
			for _, file := range files {
				e := v.Call(s, file)
				if e != nil {
					log.With("file", file).Error(e)
				}
			}
		}
	}
}

// BeforeRun ...
func (s *Slice) BeforeRun(seed Seeder) {
	s.Seeder = seed

}

// AfterRun ...
func (s *Slice) AfterRun(seed Seeder) {
}

// GetFiles ...
func GetFiles(p string) (files []string) {
	info, e := os.Stat(p)
	if e != nil {
		return nil
	}
	if info.IsDir() {
		file, e := os.Open(p)
		if e != nil {
			return nil
		}
		defer file.Close()
		names, e := file.Readdirnames(-1)
		if e != nil {
			return nil
		}
		var fullPath string
		for _, name := range names {
			fullPath = filepath.Join(p, name)
			tmp := GetFiles(fullPath)
			if tmp != nil {
				files = append(files, tmp...)
			}
		}
		return files
	}
	return append(files, p)
}

type sliceCall struct {
	path string
	//*Slice
	skipType    []interface{}
	skipExist   bool
	skipSlice   bool
	scale       int64
	sliceOutput string
}

// Path ...
func (c *sliceCall) Path() string {
	return c.path
}

type unfinishedSlice struct {
	unfinished *model.Unfinished
	format     *cmd.StreamFormat
	file       string
	sliceCall
}

// SliceCall ...
func SliceCall(call *Slice, path string) SliceCaller {
	return &sliceCall{
		path:        path,
		skipType:    call.SkipType,
		skipExist:   call.SkipExist,
		skipSlice:   call.SkipSlice,
		scale:       call.Scale,
		sliceOutput: call.SliceOutput,
	}
}

func skip(format *cmd.StreamFormat) bool {
	video := format.Video()
	audio := format.Audio()
	if audio == nil || video == nil {
		return true
	}
	if video.CodecName != "h264" || audio.CodecName != "aac" {
		return true
	}
	return false
}

// Call ...
func (c *sliceCall) Call(s *Slice, path string) (e error) {
	u := new(unfinishedSlice)
	u.sliceCall = *c
	u.unfinished = defaultUnfinished(path)
	u.unfinished.Relate = onlyName(path)
	if isPicture(path) {
		u.unfinished.Type = model.TypePoster
	} else {
		//fix name and get format
		u.format, e = parseUnfinishedFromStreamFormat(path, u.unfinished)
		if e != nil {
			return e
		}
	}
	log.Infof("%+v", u.format)
	if SkipTypeVerify("video", c.skipType...) {
		e = s.PushTo(StepperDatabase, DatabaseCallback(u, func(database *Database, eng *xorm.Engine, v interface{}) (e error) {
			u := v.(*unfinishedSlice)
			session := eng.Where("checksum = ?", u.unfinished.Checksum).
				Where("type = ?", u.unfinished.Type)
			if !model.IsExist(session, u.unfinished) || !u.skipExist {
				e = s.PushTo(StepperAPI, APICallback(u, func(api *API, api2 *httpapi.HttpApi, v interface{}) (e error) {
					us := v.(*unfinishedSlice)
					reader, e := os.Open(us.file)
					if e != nil {
						return e
					}

					resolved, e := api2.Unixfs().Add(context.Background(), files.NewReaderFile(reader))
					if e != nil {
						return e
					}
					u.unfinished.Hash = model.PinHash(resolved)
					e = api.PushTo(StepperDatabase, DatabaseCallback(u.unfinished, func(database *Database, eng *xorm.Engine, v interface{}) (e error) {
						return model.AddOrUpdateUnfinished(eng.NewSession(), v.(*model.Unfinished))
					}))
					return
				}))
			}
			return
		}))
		if e != nil {
			return
		}
	}

	if u.unfinished.Type == model.TypeVideo && !skip(u.format) {
		u1 := new(unfinishedSlice)
		u1.unfinished = u.unfinished.Clone()
		u1.unfinished.Type = model.TypeSlice
		e = s.PushTo(StepperDatabase, DatabaseCallback(u1, func(database *Database, eng *xorm.Engine, v interface{}) (e error) {
			u := v.(*unfinishedSlice)
			session := eng.Where("checksum = ?", u.unfinished.Checksum).
				Where("type = ?", u.unfinished.Type)
			if !model.IsExist(session, u.unfinished) || !u1.skipExist {
				e = database.PushTo(StepperAPI, APICallback(u, func(api *API, api2 *httpapi.HttpApi, v interface{}) (e error) {
					var sa *cmd.SplitArgs
					sa, e = sliceVideo(u, u.format)
					if e != nil {
						return e
					}
					stat, err := os.Lstat(sa.Output)
					if err != nil {
						return err
					}

					sf, err := files.NewSerialFile(sa.Output, false, stat)
					if err != nil {
						return err
					}

					resolved, e := api2.Unixfs().Add(context.Background(), sf)
					u.unfinished.Hash = model.PinHash(resolved)
					e = api.PushTo(StepperDatabase, DatabaseCallback(u, func(database *Database, eng *xorm.Engine, v interface{}) (e error) {
						u := v.(*unfinishedSlice)
						e = model.AddOrUpdateUnfinished(eng.NewSession(), u.unfinished)
						return
					}))
					return
				}))
			}
			return
		}))
	}
	return
}

func sliceVideo(us *unfinishedSlice, format *cmd.StreamFormat) (sa *cmd.SplitArgs, e error) {
	s := us.scale
	if s != 0 {
		res := format.ResolutionInt()
		if int64(res) < s {
			s = int64(res)
		}
		sa, e = cmd.FFMpegSplitToM3U8(nil, us.path, cmd.StreamFormatOption(format), cmd.ScaleOption(s), cmd.OutputOption(us.sliceOutput))
		us.unfinished.Sharpness = fmt.Sprintf("%dP", scale(s))

	} else {
		sa, e = cmd.FFMpegSplitToM3U8(nil, us.path, cmd.StreamFormatOption(format), cmd.OutputOption(us.sliceOutput))
	}

	log.Infof("%+v", sa)
	return

}

func parseUnfinishedFromStreamFormat(file string, u *model.Unfinished) (format *cmd.StreamFormat, e error) {
	format, e = cmd.FFProbeStreamFormat(file)
	if e != nil {
		return nil, e
	}

	if format.IsVideo() {
		u.Type = model.TypeVideo
		u.Sharpness = format.Resolution()
	}
	return format, nil
}
