package seed

import "testing"

// TestUpdate ...
func TestUpdate(t *testing.T) {
	NewSeed(DatabaseOption("sqlite3", "test3.db"), Update(UpdateMethodAll, UpdateContentHash))
}
