package seed_test

import "testing"

// TestUpdate ...
func TestUpdate(t *testing.T) {
	seed := NewSeed(DatabaseOption("sqlite3", "t1.db"), Update(UpdateMethodAll, UpdateContentAll))
	seed.AfterInit(ShowSQLOption(), SyncDatabase())
	seed.Start()

	seed.Wait()

}
