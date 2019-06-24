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
}

func tmp(path string, name string) string {
	mp, e := filepath.Abs(path)
	if e != nil {
		mp, e = filepath.Abs(filepath.Dir(os.Args[0]))
		if e != nil {
			//ignore error
			mp, _ = os.UserHomeDir()
		}
	}
	return filepath.Join(mp, name)
}

// Process ...
func Process(path string) Options {
	process := &process{
		path: path,
	}
	return processOption(process)
}

func prefix(s string) (ret string) {
	ret = "/ipfs/" + s
	return
}

func (p *process) slice(unfin *model.Unfinished, format *cmd.StreamFormat, file string) (err error) {
	sa, err := cmd.FFMpegSplitToM3U8(nil, file, cmd.StreamFormatOption(format), cmd.OutputOption(p.workspace))
	if err != nil {
		return err
	}
	log.Infof("%+v", sa)
	dirs, err := rest.AddDir(sa.Output)
	if err != nil {
		return err
	}

	last := unfin.SliceObject.ParseLinks(dirs)
	if last != nil {
		unfin.Sliced = true
		unfin.SliceHash = last.Hash
	}
	return nil
}

func fixPath(file string) string {
	n := strings.Replace(file, " ", "", -1)
	dir, _ := filepath.Split(n)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		log.Error(err)
	}
	err = os.Rename(file, n)
	if err != nil {
		log.Error(err)
	}
	return n
}

// Run ...
func (p *process) Run(ctx context.Context) {
	files := p.getFiles(p.path)
	log.Info(files)
	var unfin *model.Unfinished
	for _, oldFile := range files {
		file := fixPath(oldFile)
		log.With("old", oldFile, "new", file).Info("print filename")
		select {
		case <-ctx.Done():
			if err := ctx.Err(); err != nil {
				log.Error(err)
			}
			return
		default:
			log.With("file", file).Info("process run")
			unfin = DefaultUnfinished(file)
			object, err := rest.AddFile(file)
			if err != nil {
				log.Error(err)
				continue
			}
			//fix name and get format
			format, err := parseUnfinishedFromStreamFormat(file, unfin)
			if err != nil {
				log.Error(err)
				continue
			}
			log.Infof("%+v", format)

			unfin.Hash = object.Hash
			unfin.Object.Link = model.ObjectToVideoLink(object)
			if unfin.IsVideo {
				err := p.slice(unfin, format, file)
				if err != nil {
					log.With("split", file).Error(err)
					continue
				}
				p.unfinished[unfin.SliceHash] = unfin
			}

		}
		p.unfinished[unfin.Hash] = unfin
	}
	log.Infof("unfinished:%+v", p.unfinished)
	return
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

func parseUnfinishedFromStreamFormat(file string, u *model.Unfinished) (format *cmd.StreamFormat, e error) {
	format, e = cmd.FFProbeStreamFormat(file)
	if e != nil {
		return nil, e
	}

	if format.IsVideo() {
		u.IsVideo = true
		u.Sharpness = format.Resolution()
	}
	return format, nil
}

// DefaultUnfinished ...
func DefaultUnfinished(name string) *model.Unfinished {
	_, file := filepath.Split(name)
	uncat := &model.Unfinished{
		Model:       model.Model{},
		Checksum:    "",
		Type:        "other",
		Name:        file,
		Hash:        "",
		SliceHash:   "",
		IsVideo:     false,
		Sharpness:   "",
		Sync:        false,
		Sliced:      false,
		Encrypt:     false,
		Key:         "",
		M3U8:        "media.m3u8",
		Caption:     "",
		SegmentFile: "media-%05d.ts",
		Object:      new(model.VideoObject),
		SliceObject: new(model.VideoObject),
	}
	log.With("file", name).Info("calculate checksum")
	uncat.Checksum = model.Checksum(name)
	return uncat
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
