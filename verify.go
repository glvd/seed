package seed

import (
	cmd "github.com/godcong/go-ffmpeg-cmd"
	"os"
	"path"
	"path/filepath"
)

// Verify ...
type Verify struct {
	sfs  map[string]*cmd.StreamFormat
	path string
}

// NewVerify ...
func NewVerify(path string) *Verify {
	return &Verify{
		path: path,
	}
}

func isVideo(filename string) bool {
	vlist := []string{
		".swf", "flv", ".3gp", "ogm", ".vob", ".m4v", ".mkv", ".mp4", ".mpg", ".mpeg",
		".avi", ".rm", ".rmvb", ".mov", ".wmv", ".asf", ".dat", ".asx", ".wvx", ".mpe", ".mpa",
	}
	ext := path.Ext(filename)
	for _, v := range vlist {
		if ext == v {
			return true
		}
	}
	return false
}

// Check ...
func (v *Verify) Check() (sfs map[string]*cmd.StreamFormat) {
	files := GetFiles(v.path)
	sfs = make(map[string]*cmd.StreamFormat, len(files))
	for _, f := range files {
		if !isVideo(f) {
			continue
		}
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

// FailedList ...
func (v *Verify) FailedList() (failed []string) {
	if v.sfs == nil {
		v.sfs = v.Check()
	}

	for name, vv := range v.sfs {
		if vv.ResolutionInt() < 720 {
			failed = append(failed, name)
		}
	}
	return
}
