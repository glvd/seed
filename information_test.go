package seed

import "testing"

// TestInformation ...
func TestInformation(t *testing.T) {
	seed := NewSeed(Information("D:\\videoall\\video2.json", InfoFlagBSON), DatabaseOption("sqlite3", "test.db"))
	seed.Workspace = "D:\\videoall"
	seed.AfterInit(SyncDatabase())
	seed.Start()

	seed.Wait()
}

// TestName ...
func TestName(t *testing.T) {
	t.Log(onlyNo("file-09B.name"))
	t.Log(onlyNo("file-09.name"))
	t.Log(onlyNo("file"))
	t.Log(onlyNo(".file"))
	t.Log(onlyNo("."))
	t.Log(onlyNo(""))
}
