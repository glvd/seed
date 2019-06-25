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
	t.Log(onlyName("file.name"))
	t.Log(onlyName("file"))
	t.Log(onlyName(".file"))
	t.Log(onlyName("."))
	t.Log(onlyName(""))
}
