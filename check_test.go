package seed

import "testing"

// TestCheck ...
func TestCheck(t *testing.T) {
	seed := NewSeed(DatabaseOption("sqlite3", "d:\\cs.db"), Check("recursive"))
	seed.Start()
	seed.Wait()
}
