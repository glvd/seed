package seed

import "testing"

// TestInformation ...
func TestInformation(t *testing.T) {
	seed := NewSeed(Information("D:\\videoall\\video2.json", InfoFlagBSON, InfoStatusAdd), DatabaseOption("sqlite3", "test.db"), Update())
	seed.Workspace = "D:\\videoall"

	seed.Start()

	seed.Wait()
}
