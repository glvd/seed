package task_test

import (
	"fmt"
	"testing"

	"github.com/glvd/seed"
	"github.com/glvd/seed/model"
)

// TestInformation ...
func TestInformation(t *testing.T) {

}

// TestInformation2 ...
func TestInformation2(t *testing.T) {
	//info := seed.NewInformation()
	//info.ResourcePath =
	//info.Path =
	//info.InfoType = seed.InfoTypeBSON
	inf := &Information{
		InfoType:     InfoTypeBSON,
		Path:         "D:\\ipfstest\\video.json",
		ResourcePath: "D:\\ipfstest\\",
		//ProcList:     nil,
		//Start:        0,
	}
	sdb := seed.NewDatabase(model.MustDatabase(model.InitSQLite3("test.db")))
	sdb.RegisterSync(model.Video{}, model.Pin{}, model.Unfinished{})
	//
	api := seed.NewAPI("/ip4/127.0.0.1/tcp/5001")
	proc := seed.NewProcess()
	s := seed.NewSeed(sdb, api, proc)
	//
	s.Start()
	fmt.Println("waiting end")

	s.AddTasker(inf)
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
