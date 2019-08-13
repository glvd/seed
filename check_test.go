package seed

import "testing"

// TestCheck ...
func TestCheck(t *testing.T) {
	seed := NewSeed(DatabaseOption("sqlite3", "test.db"), Check(CheckTypeArg("recursive")))
	seed.AfterInit(SyncDatabase())
	seed.Start()
	seed.Wait()
}
