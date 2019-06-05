package seed

import "testing"

// TestGetFiles ...
func TestGetFiles(t *testing.T) {
	process := NewProcess("D:\\video")
	process.Run(3)

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
