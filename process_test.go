package seed

import (
	"testing"

	"github.com/glvd/seed/model"

	_ "github.com/mattn/go-sqlite3"
)

// TestGetFiles ...
func TestGetFiles(t *testing.T) {

	unfin, _ := model.FindUnfinished(nil, "d0384c15ca0862e3d558e9d610219a4bc9433f74")

	seed := NewSeed(IgnoreOption("d:\\videos\\tmp"), Process("D:\\videoall\\videos"), UnfinishedOption(unfin), Pin(PinStatusArg(PinStatusAll)))
	seed.Workspace = "d:\\videos\\tmp"
	//seed.ProcessPath ="d:\\videos"
	seed.Start()
	seed.Wait()

}

// TestProcess ...
func TestProcess(t *testing.T) {
	seed := NewSeed(Process("D:\\videoall\\videos"), SkipConvertOption(), DatabaseOption("sqlite3", "t1.db"))
	seed.Workspace = "d:\\videos\\tmp"
	//seed.ProcessPath ="d:\\videos"
	seed.AfterInit(SyncDatabase())
	seed.Start()
	seed.Wait()
}

// TestName ...
func TestName(t *testing.T) {
	t.Log(onlyNo("file-09-B.name"))
	t.Log(onlyNo("file-09B.name"))
	t.Log(onlyNo("file-001R"))
	t.Log(onlyNo(".file"))
	t.Log(onlyNo("."))
	t.Log(onlyNo(""))
	t.Log(NumberIndex("file-09"))
	t.Log(NumberIndex("file-09-C"))
	t.Log(NumberIndex("file-09-B"))
	t.Log(NumberIndex("file-09-A"))

}
