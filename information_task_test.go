package seed_test

import (
	"fmt"
	"testing"

	"github.com/glvd/seed"
	"github.com/glvd/seed/model"
)

// TestInformation ...
func TestInformation(t *testing.T) {
	info := seed.NewInformation()
	info.ResourcePath = "D:\\videoall\\images"
	info.Path = "D:\\videoall\\video3.json"
	info.InfoType = seed.InfoTypeBSON
	s := seed.NewSeed(info)
	s.Start()

	e := seed.SplitCall(info, 10000)
	if e != nil {
		t.Error(e)
	}
	s.Wait()
}

// TestInformation2 ...
func TestInformation2(t *testing.T) {
	info := seed.NewInformation()
	info.ResourcePath = "D:\\videoall"
	info.Path = "D:\\videoall\\video3.json"
	info.InfoType = seed.InfoTypeBSON
	sdb := seed.NewDatabase(model.MustDatabase(model.InitSQLite3("test.db")))
	sdb.RegisterSync(model.Video{}, model.Pin{}, model.Unfinished{})

	api := seed.NewAPI("/ip4/127.0.0.1/tcp/5001")

	s := seed.NewSeed(info, sdb, api)

	s.Start()
	fmt.Println("waiting end")
	//e := seed.SplitCall(info, 10000)
	//if e != nil {
	//	t.Error(e)
	//}
	if err := s.PushTo(seed.InformationCall(info.InfoType, info.Path)); err != nil {
		t.Error(err)
	}
	s.Wait()
	fmt.Println("waiting db end")
	sdb.Done()
	fmt.Println("db end")

}
