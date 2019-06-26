package seed

import (
	"context"
	"crypto/md5"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
	skipConvert bool
}

// BeforeRun ...
func (p *process) BeforeRun(seed *Seed) {
	p.unfinished = seed.Unfinished
	p.workspace = seed.Workspace
	p.shell = seed.Shell
	p.ignores = seed.ignores
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

func (p *process) sliceAdd(unfin *model.Unfinished, format *cmd.StreamFormat, file string) (err error) {
	sa, err := cmd.FFMpegSplitToM3U8(nil, file, cmd.StreamFormatOption(format), cmd.OutputOption(p.workspace))
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
		unfin.Type = model.TypeSlice
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
			unfin.Relate = onlyNo(file)
			if isPicture(file) {
				unfin.Type = model.TypePoster
			}
			err := p.fileAdd(unfin, file)
			if err != nil {
				log.With("add file", file).Error(err)
				continue
			}
			p.unfinished[unfin.Hash] = unfin
			if unfin.Type == model.TypePoster {
				continue
			}
			//fix name and get format
			format, err := parseUnfinishedFromStreamFormat(file, unfin)
			if err != nil {
				log.Error(err)
				continue
			}
			log.Infof("%+v", format)
			if unfin.Type == model.TypeVideo || p.skip(format) {
				unfinSlice := cloneUnfinished(unfin)
				err := p.sliceAdd(unfinSlice, format, file)
				if err != nil {
					log.With("add slice", file).Error(err)
					continue
				}
				p.unfinished[unfinSlice.Hash] = unfinSlice
			}
		}
		p.moves[unfin.Checksum] = file
	}
	return
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
	log.Info("check ", name)
	_, b = p.ignores[PathMD5(strings.ToLower(name))]
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
		return p.skipConvert
	}
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
