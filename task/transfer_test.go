package task_test

import (
	"fmt"
	"testing"

	"github.com/glvd/seed"
	"github.com/glvd/seed/model"
	"github.com/glvd/seed/task"
)

// TestTransfer ...
func TestTransfer(t *testing.T) {
	dbt := task.NewDBTransfer(model.MustDatabase(model.InitSQLite3("cs.db")))
	//jst := task.NewJSONTransfer()
	sdb := seed.NewDatabase(model.MustDatabase(model.InitSQLite3("test.db")))
	sdb.RegisterSync(model.Video{}, model.Pin{}, model.Unfinished{})
	//
	api := seed.NewAPI("/ip4/127.0.0.1/tcp/5001")
	proc := seed.NewProcess()

	s := seed.NewSeed(sdb, api, proc)
	//
	s.Start()
	fmt.Println("waiting end")

	s.AddTasker(dbt)

	s.Wait()
}
