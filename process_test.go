package seed

import (
	"github.com/yinhevr/seed/model"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

// TestGetFiles ...
func TestGetFiles(t *testing.T) {

	unfin, _ := model.FindUnfinished(nil, "d0384c15ca0862e3d558e9d610219a4bc9433f74")

	seed := NewSeed(IgnoreOption("d:\\video\\tmp"), Process("D:\\video"), UnfinishedOption(unfin), Pin(PinFlagAll))
	seed.Workspace = "d:\\video\\tmp"
	//seed.ProcessPath ="d:\\video"
	seed.Start()
	seed.Wait()

}
