package seed

import "testing"
import _ "github.com/mattn/go-sqlite3"

// TestGetFiles ...
func TestGetFiles(t *testing.T) {
	Rest()
	process := NewProcess("D:\\video")
	process.Run()

}

// TestDefaultUnfinished ...
func TestDefaultUnfinished(t *testing.T) {
	process := NewProcess("D:\\video")
	files := process.getFiles("D:\\video")
	for _, f := range files {
		Unfinished := DefaultUnfinished(f)

		log.Infof("%+v", Unfinished)
	}

}
