package seed

import (
	"github.com/yinhevr/seed/model"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

// TestGetFiles ...
func TestGetFiles(t *testing.T) {

	unfin, _ := model.FindUnfinished(nil, "d0384c15ca0862e3d558e9d610219a4bc9433f74")

	seed := NewSeed(IgnoreOption("d:\\videos\\tmp"), Process("D:\\videos"), UnfinishedOption(unfin), Pin())
	seed.Workspace = "d:\\videos\\tmp"
	//seed.ProcessPath ="d:\\videos"
	seed.Start()
	seed.Wait()

}
