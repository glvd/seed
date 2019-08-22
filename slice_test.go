package seed_test

import (
	"fmt"
	"testing"

	"github.com/glvd/seed"
	"github.com/glvd/seed/model"
)

// TestSliceCall ...
func TestNewSlice(t *testing.T) {
	sli := seed.NewSlice()
	info := seed.NewInformation()
	info.ResourcePath = "D:\\videoall"
	info.Path = "D:\\videoall\\video4.json"
	info.InfoType = seed.InfoTypeBSON
	sdb := seed.NewDatabase(model.MustDatabase(model.InitSQLite3("test.db")))
	sdb.RegisterSync(model.Video{}, model.Pin{}, model.Unfinished{})

	api := seed.NewAPI("/ip4/127.0.0.1/tcp/5001")

	s := seed.NewSeed(info, sdb, api, sli)
	s.Start()
	fmt.Println("waiting end")
	e := s.PushTo(seed.StepperSlice, seed.SliceCall("D:\\videoall\\videos", func(s *seed.Slice, v interface{}) (e error) {
		fmt.Println("slice call")
		return
	}))
	if e != nil {
		t.Error(e)
	}
	s.Wait()
	fmt.Println("waiting db end")
	sdb.Done()
	api.Done()
	fmt.Println("db end")

}
