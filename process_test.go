package seed

import "testing"
import _ "github.com/mattn/go-sqlite3"

// TestGetFiles ...
func TestGetFiles(t *testing.T) {
	Rest()
	process := NewProcess("D:\\video")
	err := process.Run(3)
	log.Info(err)

}

// TestDefaultUncategorized ...
func TestDefaultUncategorized(t *testing.T) {
	process := NewProcess("D:\\video")
	files := process.getFiles("D:\\video")
	for _, f := range files {
		uncategorized := DefaultUncategorized(f)

		log.Infof("%+v", uncategorized)
	}

}
