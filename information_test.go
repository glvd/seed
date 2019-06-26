package seed

import "testing"

// TestInformation ...
func TestInformation(t *testing.T) {
	seed := NewSeed(Information("D:\\videoall\\vtest.json", InfoFlagBSON), DatabaseOption("sqlite3", "t1.db"))
	seed.Workspace = "D:\\videoall2"
	seed.AfterInit(SyncDatabase())
	seed.Start()

	seed.Wait()
}
