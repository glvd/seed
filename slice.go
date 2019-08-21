package seed

import (
	"context"
	"os"
	"path/filepath"

	"github.com/glvd/seed/model"
	cmd "github.com/godcong/go-ffmpeg-cmd"
)

// Slice ...
type Slice struct {
	Seeder
	Scale    int64
	Path     string
	skipType []interface{}
}

// Run ...
func (s *Slice) Run(context.Context) {
	log.Info("slice running")
	files := GetFiles(s.Path)
	var unfin *model.Unfinished
	for _, file := range files {
		unfin = defaultUnfinished(file)
		unfin.Relate = onlyName(file)
		var format *cmd.StreamFormat
		var e error
		if isPicture(file) {
			unfin.Type = model.TypePoster
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

		if SkipTypeVerify("video", s.skipType...) {
			if !model.IsExist(nil, unfin) || !s.skipExist {
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
			if !model.IsExist(nil, unfinSlice) || !p.skipExist {
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
