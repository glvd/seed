package seed_test

import "testing"

// TestPin ...
func TestPin(t *testing.T) {
	seed := NewSeed(DatabaseOption("sqlite3", "data.db"), Pin(PinStatusArg(PinStatusAll)))
	//seed.Workspace = "D:\\videoall"
	//seed.AfterInit()
	seed.Start()
	seed.Wait()
}
