package seed

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/glvd/seed/model"
	cmd "github.com/godcong/go-ffmpeg-cmd"
)

// SliceCaller ...
type SliceCaller interface {
	Call(*Slice) error
}

// Slice ...
type Slice struct {
	*Thread
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

// Option ...
func (s *Slice) Option(seed Seeder) {
	sliceOption(s)(seed)
}

func sliceOption(slice *Slice) Options {
	return func(seeder Seeder) {
		seeder.SetBaseThread(StepperSlice, slice)
	}
}

// Push ...
func (s *Slice) push(cb interface{}) error {
	if v, b := cb.(SliceCaller); b {
		s.cb <- v
		return nil
	}
	return errors.New("not slice callback")
}

// NewSlice ...
func NewSlice() *Slice {
	return &Slice{
		cb:     make(chan SliceCaller),
		Thread: NewThread(),
	}
}

// Run ...
func (s *Slice) Run(ctx context.Context) {
	log.Info("slice running")
SliceEnd:
	for {
		select {
		case <-ctx.Done():
			break SliceEnd
		case v := <-s.cb:
			if v == nil {
				break SliceEnd
			}
			s.SetState(StateRunning)
			e := v.Call(s)
			if e != nil {
				log.Error(e)
			}
		case <-time.After(30 * time.Second):
			log.Info("slice time out")
			s.SetState(StateWaiting)
		}
	}
	close(s.cb)
	s.Finished()
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
	file       string
}

// SliceCallbackFunc ...
type SliceCallbackFunc func(s *Slice, sa *cmd.SplitArgs, v interface{}) (e error)

// SliceCall ...
func SliceCall(file string, u *model.Unfinished, cb SliceCallbackFunc) (Stepper, SliceCaller) {
	return StepperSlice, &sliceCall{
		cb:         cb,
		unfinished: u,
		file:       file,
	}
}

func scale(scale int64) int {
	switch scale {
	case 480, 1080:
		return int(scale)
	default:
		return 720
	}
}

func scaleStr(s int64) string {
	return fmt.Sprintf("%dP", scale(s))
}

func skip(format *cmd.StreamFormat) bool {
	video := format.Video()
	audio := format.Audio()
	if audio == nil || video == nil {
		return true
	}
	return false
}

// Call ...
func (c *sliceCall) Call(s *Slice) (e error) {
	sa, e := sliceVideo(s, c.file, c.unfinished)
	if e != nil {
		return e
	}
	return c.cb(s, sa, c.unfinished)
}

func sliceVideo(slice *Slice, file string, u *model.Unfinished) (sa *cmd.SplitArgs, e error) {
	format, e := cmd.FFProbeStreamFormat(file)
	if e != nil {
		return nil, e
	}
	if skip(format) {
		return nil, errors.New("format video/audio not found")
	}

	u.Type = model.TypeSlice
	s := slice.Scale
	if s != 0 {
		res := format.ResolutionInt()
		if int64(res) < s {
			s = int64(res)
		}
		sa, e = cmd.FFMpegSplitToM3U8(nil, file, cmd.StreamFormatOption(format), cmd.ScaleOption(s), cmd.OutputOption(slice.SliceOutput))
		u.Sharpness = scaleStr(s)
	} else {
		sa, e = cmd.FFMpegSplitToM3U8(nil, file, cmd.StreamFormatOption(format), cmd.OutputOption(slice.SliceOutput))
		u.Sharpness = format.Resolution() + "P"
	}

	log.Infof("%+v", sa)
	return
}

type sliceCall struct {
	cb         SliceCallbackFunc
	unfinished *model.Unfinished
	file       string
}

var _ SliceCaller = &sliceCall{}
