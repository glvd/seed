package seed

import "testing"

// TestInformation ...
func TestInformation(t *testing.T) {
	seed := NewSeed(Information("D:\\videoall\\video2.json", InfoFlagBSON), DatabaseOption("sqlite3", "test2.db"), Update(UpdateStatusAdd))
	seed.Workspace = "D:\\videoall"
	seed.AfterInit(SyncDatabase())
	seed.Start()

	seed.Wait()
}