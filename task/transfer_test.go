package task_test

import (
	"fmt"
	"testing"

	"github.com/glvd/seed"
	"github.com/glvd/seed/model"
	"github.com/glvd/seed/task"
)

func TestTransferDB(t *testing.T) {
	dbt := task.NewDBTransfer(model.MustDatabase(model.InitSQLite3("cs.db")))
	sdb := seed.NewDatabase(model.MustDatabase(model.InitSQLite3("test.db")))
	sdb.RegisterSync(model.Video{}, model.Pin{}, model.Unfinished{})
	api := seed.NewAPI("/ip4/127.0.0.1/tcp/5001")
	proc := seed.NewProcess()

	s := seed.NewSeed(sdb, api, proc)
	s.Start()

	s.AddTasker(dbt)
	fmt.Println("waiting end")
	s.Wait()
}

// TestTransfer ...
func TestTransferJSON(t *testing.T) {
	//dbt := task.NewDBTransfer(model.MustDatabase(model.InitSQLite3("0916.db")))
	jst := task.NewJSONTransfer("output.json")
	//jst := task.NewJSONTransfer()
	//dbt.Status = task.TransferStatusToJSON
	sdb := seed.NewDatabase(model.MustDatabase(model.InitSQLite3("0916.db")))
	sdb.RegisterSync(model.Video{}, model.Pin{}, model.Unfinished{})
	//
	//api := seed.NewAPI("/ip4/127.0.0.1/tcp/5001")
	//proc := seed.NewProcess()

	s := seed.NewSeed(sdb)
	//
	s.Start()

	s.AddTasker(jst)
	fmt.Println("waiting end")
	s.Wait()
}
