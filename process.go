package seed

import (
	"context"
	"crypto/md5"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	cmd "github.com/godcong/go-ffmpeg-cmd"

	shell "github.com/godcong/go-ipfs-restapi"
	"github.com/yinhevr/seed/model"
)

func dummy(process *process) (e error) {
	log.Info("dummy called")
	return
}

// process ...
type process struct {
	workspace   string
	path        string
	shell       *shell.Shell
	ignores     map[string][]byte
	unfinished  map[string]*model.Unfinished
	moves       map[string]string
	scale       int64
	skipConvert bool
	skipSource  bool
	skipExist   bool
	noSlice     bool
	preAdd      bool
}

// BeforeRun ...
func (p *process) BeforeRun(seed *Seed) {
	p.unfinished = seed.Unfinished
	p.workspace = seed.Workspace
	p.shell = seed.Shell
	p.moves = seed.Moves
	p.preAdd = seed.preAdd
	p.scale = seed.Scale
	p.skipConvert = seed.skipConvert
	p.skipSource = seed.skipSource
	p.ignores = seed.ignores
	p.skipExist = seed.skipExist

}

// AfterRun ...
func (p *process) AfterRun(seed *Seed) {
	seed.Unfinished = p.unfinished
	seed.Moves = p.moves
}

// Process ...
func Process(path string) Options {
	process := &process{
		path: path,
	}
	return processOption(process)
}

func scale(scale int64) int {
	switch scale {
	case 480, 1080:
		return int(scale)
	default:
		return 720
	}
}

func (p *process) sliceAdd(unfin *model.Unfinished, format *cmd.StreamFormat, file string) (err error) {
	var sa *cmd.SplitArgs
	s := p.scale
	if s != 0 {
		res := format.ResolutionInt()
		if int64(res) < s {
			s = int64(res)
		}
		sa, err = cmd.FFMpegSplitToM3U8(nil, file, cmd.StreamFormatOption(format), cmd.ScaleOption(s), cmd.OutputOption(p.workspace))
		unfin.Sharpness = fmt.Sprintf("%dP", scale(s))
	} else {
		sa, err = cmd.FFMpegSplitToM3U8(nil, file, cmd.StreamFormatOption(format), cmd.OutputOption(p.workspace))
	}

	if err != nil {
		return err
	}
	log.Infof("%+v", sa)
	dirs, err := p.shell.AddDir(sa.Output)
	if err != nil {
		return err
	}

	last := unfin.Object.ParseLinks(dirs)
	if last != nil {
		unfin.Hash = last.Hash
	}
	return model.AddOrUpdateUnfinished(unfin)
}

func (p *process) fileAdd(unfin *model.Unfinished, file string) (err error) {
	object, err := p.shell.AddFile(file)
	if err != nil {
		log.Error(err)
		return
	}
	unfin.Hash = object.Hash
	unfin.Object.Link = model.ObjectToVideoLink(object)
	return model.AddOrUpdateUnfinished(unfin)
}

func onlyName(name string) string {
	_, name = filepath.Split(name)
	for i := len(name) - 1; i >= 0 && !os.IsPathSeparator(name[i]); i-- {
		if name[i] == '.' {
			return name[:i]
		}
	}
	return ""
}

func onlyNo(name string) string {
	s := []rune(onlyName(name))
	last := len(s) - 1
	if last > 0 && unicode.IsLetter(s[last]) {
		return string(s[:last])
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
		return strings.IndexByte(relateList, name[size-1])
	}
	return -1
}

// Run ...
func (p *process) Run(ctx context.Context) {
	files := p.getFiles(p.path)
	log.Info(files)
	var unfin *model.Unfinished
	for _, file := range files {
		select {
		case <-ctx.Done():
			if err := ctx.Err(); err != nil {
				log.Error(err)
			}
			return
		default:
			log.With("file", file).Info("process run")
			unfin = defaultUnfinished(file)
			unfin.Relate = onlyName(file)
			var format *cmd.StreamFormat
			var e error
			if isPicture(file) {
				unfin.Type = model.TypePoster
				p.unfinished[unfin.Hash] = unfin
			} else {
				//fix name and get format
				format, e = parseUnfinishedFromStreamFormat(file, unfin)
				if e != nil {
					log.Error(e)
					continue
				}
			}
			log.Infof("%+v", format)
			log.Info("adding:", file)

			if p.skipSource {
				if !unfin.IsExist() || !p.skipExist {
					err := p.fileAdd(unfin, file)
					if err != nil {
						log.With("add file", file).Error(err)
						continue
					}
					p.unfinished[unfin.Hash] = unfin
				}

			}

			if unfin.Type == model.TypeVideo && !p.skip(format) {
				unfinSlice := unfin.Clone()
				unfinSlice.Type = model.TypeSlice
				if !unfin.IsExist() || !p.skipExist {
					if p.noSlice {
						continue
					}
					err := p.sliceAdd(unfinSlice, format, file)
					if err != nil {
						log.With("add slice", file).Error(err)
						continue
					}
					p.unfinished[unfinSlice.Hash] = unfinSlice
				}

			}
		}

		p.moves[file] = unfin.Hash
	}
	return
}

func updateRelate(unfinished model.Unfinished) error {
	return nil
}

func isPicture(name string) bool {
	picture := ".bmp,.jpg,.png,.tif,.gif,.pcx,.tga,.exif,.fpx,.svg,.psd,.cdr,.pcd,.dxf,.ufo,.eps,.ai,.raw,.WMF,.webp"
	ext := filepath.Ext(name)
	return strings.Index(picture, ext) != -1
}

// PathMD5 ...
func PathMD5(s ...string) string {
	str := filepath.Join(s...)
	return fmt.Sprintf("%x", md5.Sum([]byte(str)))
}

// CheckIgnore ...
func (p *process) CheckIgnore(name string) (b bool) {
	if p.ignores == nil {
		return false
	}
	log.Info("noCheck ", name)
	_, b = p.ignores[name]
	return
}

func (p *process) getFiles(ws string) (files []string) {
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
			if p.CheckIgnore(fullPath) {
				continue
			}
			tmp := p.getFiles(fullPath)
			if tmp != nil {
				files = append(files, tmp...)
			}
		}
		return files
	}
	return append(files, ws)
}

func (p *process) skip(format *cmd.StreamFormat) bool {
	if !p.skipConvert {
		log.Info("noskip")
		return p.skipConvert
	}
	video := format.Video()
	audio := format.Audio()
	if audio == nil || video == nil {
		log.Info("skip")
		return true
	}
	if video.CodecName != "h264" || audio.CodecName != "aac" {
		log.Info("skip")
		return true
	}
	log.Info("noskip")
	return false
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

// Load ...
func Load(path string) []*VideoSource {
	var vs []*VideoSource
	e := ReadJSON(path, &vs)
	if e != nil {
		return nil
	}
	return vs
}
