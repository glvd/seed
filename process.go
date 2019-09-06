package seed

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"

	"github.com/glvd/seed/model"
	cmd "github.com/godcong/go-ffmpeg-cmd"
)

func dummy(process *Process) (e error) {
	log.Info("dummy called")
	return
}

// Process ...
type Process struct {
	*Thread
	//taskMutex *sync.RWMutex
	//tasks     map[string]*Task
	cb chan ProcessCaller
	//workspace   string
	//path        string
	//moves       map[string]string
	//scale int64
	//skipConvert bool
	//skipExist   bool
	//noSlice     bool
	//preAdd      bool
	//skipType    []interface {}
}

// Push ...
func (p *Process) Push(v interface{}) error {
	return p.push(v)
}

func (p *Process) push(cb interface{}) error {
	if v, b := cb.(ProcessCaller); b {
		p.cb <- v
		return nil
	}
	return errors.New("not process callback")

}

// NewProcess ...
func NewProcess() *Process {
	process := &Process{}
	process.Thread = NewThread()
	process.cb = make(chan ProcessCaller)
	return process
}

// Option ...
func (p *Process) Option(seeder Seeder) {
	processOption(p)(seeder)
}

func (p *Process) sliceAdd(unfin *model.Unfinished, format *cmd.StreamFormat, file string) (err error) {
	//var sa *cmd.SplitArgs
	//s := int64(0) // p.scale
	//if s != 0 {
	//	res := format.ResolutionInt()
	//	if int64(res) < s {
	//		s = int64(res)
	//	}
	//	sa, err = cmd.FFMpegSplitToM3U8(nil, file, cmd.StreamFormatOption(format), cmd.ScaleOption(s), cmd.OutputOption(p.workspace))
	//
	//} else {
	//	sa, err = cmd.FFMpegSplitToM3U8(nil, file, cmd.StreamFormatOption(format), cmd.OutputOption(p.workspace))
	//}
	//
	//if err != nil {
	//	return err
	//}
	//log.Infof("%+v", sa)
	//
	//dirs, err := p.shell.AddDir(sa.Output)
	//if err != nil {
	//	return err
	//}
	//
	//last := unfin.Object.ParseLinks(dirs)
	//if last != nil {
	//	unfin.Hash = last.Hash
	//}
	return model.AddOrUpdateUnfinished(nil, unfin)
}

func (p *Process) fileAdd(unfin *model.Unfinished, file string) (err error) {
	//object, err := p.shell.AddFile(file)
	//if err != nil {
	//	log.Error(err)
	//	return
	//}
	//unfin.Hash = object.Hash
	//unfin.Object.Link = model.ObjectToVideoLink(object)
	return model.AddOrUpdateUnfinished(nil, unfin)
}

// OnlyName ...
func OnlyName(name string) string {
	_, name = filepath.Split(name)
	for i := len(name) - 1; i >= 0 && !os.IsPathSeparator(name[i]); i-- {
		if name[i] == '.' {
			return name[:i]
		}
	}
	return ""
}

// OnlyNo ...
func OnlyNo(name string) string {
	s := []rune(OnlyName(name))
	last := len(s) - 1
	if last > 0 && unicode.IsLetter(s[last]) {
		if s[last-1] == rune('-') {
			return string(s[:last-1])
		}
		//return string(s[:last])
	}
	return string(s)
}

// RelateList ...
const relateList = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

// IndexNumber ...
func IndexNumber(index int) byte {
	if index > len(relateList) {
		return relateList[0]
	}
	return relateList[index]
}

// NumberIndex ...
func NumberIndex(name string) int {
	size := len(name)
	if size > 0 {
		return strings.Index(relateList, LastSlice(name, "@"))
	}
	return -1
}

// LastSlice ...
func LastSlice(s, sep string) string {
	ss := strings.Split(s, sep)
	for i := len(ss) - 1; i >= 0; i-- {
		if ss[i] == "" {
			continue
		}
		return ss[i]
	}
	return ""
}

// Run ...
func (p *Process) Run(ctx context.Context) {
	log.Info("process running")
ProcessEnd:
	for {
		select {
		case <-ctx.Done():
			break ProcessEnd
		case v := <-p.cb:
			if v == nil {
				break ProcessEnd
			}
			p.SetState(StateRunning)
			e := v.Call(p)
			if e != nil {
				log.Error(e)
			}
		case <-time.After(30 * time.Second):
			log.Info("process time out")
			p.SetState(StateWaiting)
		}
	}
	close(p.cb)
	p.Finished()
}

// PathMD5 ...
func PathMD5(s ...string) string {
	str := filepath.Join(s...)
	return fmt.Sprintf("%x", md5.Sum([]byte(str)))
}

func (p *Process) getFiles(ws string) (files []string) {
	info, e := os.Stat(ws)
	if e != nil {
		return nil
	}
	if info.IsDir() {
		file, e := os.Open(ws)
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
			fullPath = filepath.Join(ws, name)
			//if p.CheckIgnore(fullPath) {
			//	continue
			//}
			tmp := p.getFiles(fullPath)
			if tmp != nil {
				files = append(files, tmp...)
			}
		}
		return files
	}
	return append(files, ws)
}

func (p *Process) skip(format *cmd.StreamFormat) bool {
	//if !p.skipConvert {
	//	log.Info("noskip")
	//	return p.skipConvert
	//}
	video := format.Video()
	audio := format.Audio()
	if audio == nil || video == nil {
		log.Info("isMedia")
		return true
	}
	if video.CodecName != "h264" || audio.CodecName != "aac" {
		log.Info("isMedia")
		return true
	}
	log.Info("noskip")
	return false
}

func moveSuccess(file string) (e error) {
	dir, name := filepath.Split(file)
	newPath := filepath.Join(dir, "success")
	_ = os.MkdirAll(newPath, os.ModePerm)
	newPathFile := filepath.Join(newPath, name)
	return os.Rename(file, newPathFile)
}

// MustString  must string
func MustString(val, src string) string {
	if val != "" {
		return val
	}
	return src
}

// processOption ...
func processOption(process *Process) Options {
	return func(seed Seeder) {
		seed.SetBaseThread(StepperProcess, process)
	}
}
