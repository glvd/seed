package seed

import "testing"

// TestInformation ...
func TestInformation(t *testing.T) {
	seed := NewSeed(Information("D:\\videoall\\video.json", InfoFlagBSON), DatabaseOption("sqlite3", "t1.db"))
	seed.Workspace = "D:\\videoall"
	seed.AfterInit(SyncDatabase())
	seed.Start()

	seed.Wait()
}

// TestInformation2 ...
func TestInformation2(t *testing.T) {
	seed := NewSeed(Information("D:\\videoall\\video.json", InfoFlagBSON), DatabaseOption("sqlite3", "t1.db"), MoveInfo("D:\\videoall\\picSuccess"))
	seed.Workspace = "D:\\videoall"
	seed.AfterInit(SyncDatabase())
	seed.MaxLimit = 2
	seed.Start()

	seed.Wait()
}
