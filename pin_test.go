package seed

import "testing"

// TestPin ...
func TestPin(t *testing.T) {
	seed := NewSeed(Information("D:\\videoall\\video2.json", InfoFlagBSON), Update(UpdateStatusAdd), Pin())
	seed.Workspace = "D:\\videoall"
	seed.AfterInit(SyncDatabase(), ShowSQLOption(), ShowExecTimeOption())
	seed.Start()

	seed.Wait()
}
