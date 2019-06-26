package seed

import "testing"

// TestPin ...
func TestPin(t *testing.T) {
	seed := NewSeed(DatabaseOption("sqlite3", "t1.db"), Pin(PinStatusAll))
	//seed.Workspace = "D:\\videoall"
	//seed.AfterInit()
	seed.Start()

	seed.Wait()
}
