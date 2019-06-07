package seed

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

// TestGetFiles ...
func TestGetFiles(t *testing.T) {
	Rest()
	process := NewProcessSeeder("D:\\video")
	process.Start()
	process.Wait()
}
