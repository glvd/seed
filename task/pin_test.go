package task_test

import (
	"fmt"
	"testing"

	"github.com/glvd/seed"
	"github.com/glvd/seed/model"
	"github.com/glvd/seed/task"
)

// TestPin ...
func TestPin(t *testing.T) {
	//jst := task.NewJSONTransfer()
	p := task.NewPin()
	sdb := seed.NewDatabase(model.MustDatabase(model.InitSQLite3("test.db")))
	sdb.RegisterSync(model.Video{}, model.Pin{}, model.Unfinished{})
	//
	api := seed.NewAPI("/ip4/127.0.0.1/tcp/5001")
	proc := seed.NewProcess()

	s := seed.NewSeed(sdb, api, proc)
	//
	s.Start()
	fmt.Println("waiting end")

	s.AddTasker(p)

	s.Wait()

}
