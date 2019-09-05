package task

import (
	"fmt"
	"testing"

	"github.com/glvd/seed"
	"github.com/glvd/seed/model"
)

// TestNewVideoProcess ...
func TestNewVideoProcess(t *testing.T) {
	process := NewVideoProcess()
	process.Scale = 720
	process.Skip = []interface{}{"video"}
	process.Path = "D:\\video\\test"
	sdb := seed.NewDatabase(model.MustDatabase(model.InitSQLite3("test.db")))
	sdb.RegisterSync(model.Video{}, model.Pin{}, model.Unfinished{})
	//
	api := seed.NewAPI("/ip4/127.0.0.1/tcp/5001")
	proc := seed.NewProcess()

	slice := seed.NewSlice()
	s := seed.NewSeed(sdb, api, proc, slice)
	//
	s.Start()
	fmt.Println("waiting end")

	s.AddTasker(process)
	////e := seed.SplitCall(info, 10000)
	////if e != nil {
	////	t.Error(e)
	////}
	//if err := s.PushTo(seed.InformationCall(info.InfoType, info.Path)); err != nil {
	//	t.Error(err)
	//}
	s.Wait()
	//fmt.Println("waiting db end")
	//sdb.Done()
	//fmt.Println("db end")
}
