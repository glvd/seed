package seed

import (
	cmd "github.com/godcong/go-ffmpeg-cmd"
	"os"
	"path/filepath"
)

type Verify struct {
	path string
}

func (v *Verify) getFiles(ws string) (files []string) {
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
			tmp := v.getFiles(fullPath)
			if tmp != nil {
				files = append(files, tmp...)
			}
		}
		return files
	}
	return append(files, ws)
}

func NewVerify(path string) *Verify {
	return &Verify{
		path: path,
	}
}

func (v *Verify) Check() (sfs map[string]*cmd.StreamFormat) {
	files := v.getFiles(v.path)
	sfs = make(map[string]*cmd.StreamFormat, len(files))
	for _, f := range files {
		format, e := cmd.FFProbeStreamFormat(f)
		if e != nil {
			log.Error(e)
			continue
		}
		_, file := filepath.Split(f)
		sfs[file] = format
	}
	return sfs
}
