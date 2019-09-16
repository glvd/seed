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

// Scale ...
type Scale int64

// HighScale ...
const HighScale Scale = 1080

// MiddleScale ...
const MiddleScale Scale = 720

// LowScale ...
const LowScale Scale = 480

// SliceCaller ...
type SliceCaller interface {
	Call(*Slice) error
}

// Slice ...
type Slice struct {
	*Thread
	Scale       Scale
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
	output := os.TempDir()
	return &Slice{
		SliceOutput: output,
		cb:          make(chan SliceCaller),
		Thread:      NewThread(),
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
		case <-time.After(TimeOutLimit):
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

func scale(scale Scale) int {
	switch scale {
	case 480, 1080:
		return int(scale)
	default:
		return 720
	}
}

func scaleStr(s Scale) string {
	return fmt.Sprintf("%dP", scale(s))
}

func isMedia(format *cmd.StreamFormat) bool {
	video := format.Video()
	audio := format.Audio()
	if audio == nil || video == nil {
		return false
	}
	return true
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
	if !isMedia(format) {
		return nil, errors.New("format video/audio not found")
	}

	u.Type = model.TypeSlice
	s := slice.Scale
	if s != 0 {
		res := format.ResolutionInt()
		if int64(res) < int64(s) {
			s = Scale(res)
		}
		sa, e = cmd.FFMpegSplitToM3U8(nil, file, cmd.StreamFormatOption(format), cmd.ScaleOption(int64(s)), cmd.OutputOption(slice.SliceOutput))
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
