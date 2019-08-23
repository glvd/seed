package seed

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

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
	Threader
	Scale       int64
	SliceOutput string
	SkipType    []interface{}
	SkipExist   bool
	SkipSlice   bool
	cb          chan SliceCaller
}

// Push ...
func (s *Slice) Push(v interface{}) error {
	return s.push(v)
}

// Done ...
func (s *Slice) Done() <-chan bool {
	go func() {
		s.cb <- nil
	}()
	return s.Threader.Done()
}

// Option ...
func (s *Slice) Option(seed Seeder) {
	sliceOption(s)(seed)
}

func sliceOption(slice *Slice) Options {
	return func(seeder Seeder) {
		seeder.SetThread(StepperSlice, slice)
	}
}

// Push ...
func (s *Slice) push(cb interface{}) error {
	if v, b := cb.(SliceCaller); b {
		s.cb <- v
		return nil
	}
	return xerrors.New("not slice callback")
}

// NewSlice ...
func NewSlice() *Slice {
	return &Slice{
		cb:       make(chan SliceCaller),
		Threader: NewThread(),
	}
}

// Run ...
func (s *Slice) Run(ctx context.Context) {
	log.Info("slice running")
SliceEnd:
	for {
		select {
		case <-ctx.Done():
			s.SetState(StateStop)
			break SliceEnd
		case v := <-s.cb:
			if v == nil {
				s.SetState(StateStop)
				break SliceEnd
			}
			s.SetState(StateRunning)
			files := GetFiles(v.Path())
			for _, file := range files {
				e := v.Call(s, file)
				if e != nil {
					log.With("file", file).Error(e)
				}
			}
		case <-time.After(30 * time.Second):
			s.SetState(StateWaiting)
		}
	}
	close(s.cb)
	s.Finished()
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

type unfinishedSlice struct {
	unfinished *model.Unfinished
	format     *cmd.StreamFormat
	file       string
	sliceCall
}

// SliceCallbackFunc ...
type SliceCallbackFunc func(s *Slice, v interface{}) (e error)

// SliceOptions ...
type SliceOptions func(call *sliceCall)

// SliceOutputOption ...
func SliceOutputOption(path string) SliceOptions {
	return func(call *sliceCall) {
		call.sliceOutput = path
	}
}

// SliceCall ...
func SliceCall(path string, cb SliceCallbackFunc, sliceOption ...SliceOptions) (Stepper, SliceCaller) {
	return StepperSlice, &sliceCall{
		cb:   cb,
		path: path,
		//skipType:    call.SkipType,
		//skipExist:   call.SkipExist,
		//skipSlice:   call.SkipSlice,
		//scale:       call.Scale,
		//sliceOutput: call.SliceOutput,
	}
}

func skip(format *cmd.StreamFormat) bool {
	video := format.Video()
	audio := format.Audio()
	if audio == nil || video == nil {
		return true
	}
	//if video.CodecName != "h264" || audio.CodecName != "aac" {
	//	return true
	//}
	return false
}

// Call ...
func (c *sliceCall) Call(s *Slice, path string) (e error) {

	return c.cb(s, path)
	//u := new(unfinishedSlice)
	//u.file = path
	//u.sliceCall = *c
	//u.unfinished = defaultUnfinished(path)
	//u.unfinished.Relate = onlyName(path)
	//if isPicture(path) {
	//	u.unfinished.Type = model.TypePoster
	//} else {
	//	//fix name and get format
	//	u.format, e = parseUnfinishedFromStreamFormat(path, u.unfinished)
	//	if e != nil {
	//		return e
	//	}
	//}
	//log.Infof("%+v", u.format)
	//if !SkipTypeVerify("video", c.skipType...) {
	//	e = s.PushTo(DatabaseCallback(u, func(database *Database, eng *xorm.Engine, v interface{}) (e error) {
	//		u := v.(*unfinishedSlice)
	//		session := eng.Where("checksum = ?", u.unfinished.Checksum).
	//			Where("type = ?", u.unfinished.Type)
	//		if !model.IsExist(session, model.Unfinished{}) || !u.skipExist {
	//			log.With("file", u.file).Info("video")
	//			e = s.PushTo(APICallback(u, func(api *API, api2 *httpapi.HttpApi, v interface{}) (e error) {
	//				us := v.(*unfinishedSlice)
	//				file, err := os.Open(us.file)
	//				if err != nil {
	//					return err
	//				}
	//
	//				resolved, err := api2.Unixfs().Add(context.Background(),
	//					files.NewReaderFile(file),
	//					func(settings *options.UnixfsAddSettings) error {
	//						settings.Pin = true
	//						return nil
	//					})
	//				if err != nil {
	//					return err
	//				}
	//				u.unfinished.Hash = model.PinHash(resolved)
	//				e = api.PushTo(DatabaseCallback(u.unfinished, func(database *Database, eng *xorm.Engine, v interface{}) (e error) {
	//					u := v.(*model.Unfinished)
	//					log.With("hash", u.Hash, "relate", u.Relate).Info("update unfinished")
	//					return model.AddOrUpdateUnfinished(eng.NewSession(), v.(*model.Unfinished))
	//				}))
	//				return
	//			}))
	//		}
	//		return
	//	}))
	//	if e != nil {
	//		log.Error(e)
	//	}
	//}
	//log.With("type", u.unfinished.Type).Info("video info")
	//if u.unfinished.Type == model.TypeVideo /*&& !skip(u.format) */ {
	//	u1 := new(unfinishedSlice)
	//	u1.file = path
	//	u1.sliceCall = *c
	//	u1.unfinished = u.unfinished.Clone()
	//	u1.unfinished.Type = model.TypeSlice
	//	log.Info("slice run")
	//	e = s.PushTo(DatabaseCallback(u1, func(database *Database, eng *xorm.Engine, v interface{}) (e error) {
	//		u := v.(*unfinishedSlice)
	//		session := eng.Where("checksum = ?", u.unfinished.Checksum).
	//			Where("type = ?", u.unfinished.Type)
	//		if !model.IsExist(session, model.Unfinished{}) || !u1.skipExist {
	//			log.With("file", u.file).Info("slice")
	//			e = database.PushTo(APICallback(u, func(api *API, api2 *httpapi.HttpApi, v interface{}) (e error) {
	//				var sa *cmd.SplitArgs
	//				sa, e = sliceVideo(u, u.format)
	//				if e != nil {
	//					return e
	//				}
	//				stat, err := os.Lstat(sa.Output)
	//				if err != nil {
	//					return err
	//				}
	//
	//				sf, err := files.NewSerialFile(sa.Output, false, stat)
	//				if err != nil {
	//					return err
	//				}
	//				resolved, err := api2.Unixfs().Add(context.Background(), sf, func(settings *options.UnixfsAddSettings) error {
	//					settings.Pin = true
	//					return nil
	//				})
	//				if err != nil {
	//					return err
	//				}
	//				u.unfinished.Hash = model.PinHash(resolved)
	//				e = api.PushTo(DatabaseCallback(u, func(database *Database, eng *xorm.Engine, v interface{}) (e error) {
	//					u := v.(*unfinishedSlice)
	//					e = model.AddOrUpdateUnfinished(eng.NewSession(), u.unfinished)
	//					return
	//				}))
	//				return
	//			}))
	//		}
	//		return
	//	}))
	//}
	//return
}

func sliceVideo(us *unfinishedSlice, format *cmd.StreamFormat) (sa *cmd.SplitArgs, e error) {
	s := us.scale
	if s != 0 {
		res := format.ResolutionInt()
		if int64(res) < s {
			s = int64(res)
		}
		sa, e = cmd.FFMpegSplitToM3U8(nil, us.file, cmd.StreamFormatOption(format), cmd.ScaleOption(s), cmd.OutputOption(us.sliceOutput))
		us.unfinished.Sharpness = fmt.Sprintf("%dP", scale(s))

	} else {
		sa, e = cmd.FFMpegSplitToM3U8(nil, us.file, cmd.StreamFormatOption(format), cmd.OutputOption(us.sliceOutput))
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

type sliceCall struct {
	cb   SliceCallbackFunc
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

var _ SliceCaller = &sliceCall{}
