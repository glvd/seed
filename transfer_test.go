package seed

import "testing"

// TestTransfer ...
func TestTransfer(t *testing.T) {
	seed := NewSeed(Transfer("D:\\videoall\\video2.json", TransferFlagJSON, TransferFlagSQLite, TransferStatusUpdate), Pin(PinFlagAll), Update())
	seed.Workspace = "D:\\videoall"
	seed.Start()
	seed.Wait()
}
