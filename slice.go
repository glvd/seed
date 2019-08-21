package seed

import (
	"context"
	"os"
	"path/filepath"
)

// Slice ...
type Slice struct {
	Seeder
	Scale int64
	Path  string
}

// Run ...
func (s *Slice) Run(context.Context) {
	log.Info("slice running")
	files := GetFiles(s.Path)

	for _, file := range files {

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
