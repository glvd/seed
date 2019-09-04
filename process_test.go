package seed_test

import (
	"testing"

	"github.com/glvd/seed"
	"github.com/glvd/seed/task"
	_ "github.com/mattn/go-sqlite3"
)

// TestProcess ...
func TestProcess(t *testing.T) {
	seeder := seed.NewSeed()
	proc := seed.NewProcess()
	//tsk := seed.NewTask()
	seeder.Register(proc)

	seeder.Start()
	info := &task.Information{
		InfoType:     task.InfoTypeBSON,
		Path:         "D:\\videoall\\videos",
		ResourcePath: "",
		ProcList:     nil,
		Start:        0,
	}
	seeder.RunTask(info.Task())
	//tsk.AddTask(info)
	//proc.AddTask(seed.InformationTask(info))seed.InformationTask(info)
	seeder.Wait()

}

// TestName ...
func TestName(t *testing.T) {
	//t.Log(onlyNo("file-09-B.name"))
	//t.Log(onlyNo("file-09B.name"))
	//t.Log(onlyNo("file-001R"))
	//t.Log(onlyNo(".file"))
	//t.Log(onlyNo("."))
	//t.Log(onlyNo(""))
	//t.Log(NumberIndex("file-09"))
	//t.Log(NumberIndex("file-09-C"))
	//t.Log(NumberIndex("file-09-B"))
	//t.Log(NumberIndex("file-09-A"))
}
