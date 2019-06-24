package seed

import "testing"

// TestPin ...
func TestPin(t *testing.T) {
	seed := NewSeed(Information("D:\\videoall\\video2.json", InfoFlagBSON), Pin())
	seed.Workspace = "D:\\videoall"
	//seed.AfterInit()
	seed.Start()

	seed.Wait()
}
