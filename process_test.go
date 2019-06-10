package seed

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

// TestGetFiles ...
func TestGetFiles(t *testing.T) {
	Rest()
	//process := NewProcessSeeder("D:\\video")
	//process.Start()
	//process.Wait()
	//
	seed := NewSeed(Process("D:\\video"), Pin(PinFlagAll))
	seed.Workspace = "d:\\video"
	//seed.ProcessPath ="d:\\video"
	seed.Start()
	seed.Wait()

}
