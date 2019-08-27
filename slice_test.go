package seed_test

import (
	"fmt"
	"testing"

	"github.com/glvd/seed"
	"github.com/glvd/seed/model"
	cmd "github.com/godcong/go-ffmpeg-cmd"
)

// TestSliceCall ...
func TestNewSlice(t *testing.T) {
	sli := seed.NewSlice()
	sli.SliceOutput = "d:\\temp"

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
	e := s.PushTo(seed.SliceCall("D:\\videoall\\videos\\MIAA-087.mp4", &model.Unfinished{}, func(s *seed.Slice, sa *cmd.SplitArgs, v interface{}) (e error) {
		fmt.Println(*sa, v)
		return nil
	}))
	if e != nil {
		t.Error(e)
	}
	s.Wait()
	fmt.Println("waiting db end")
	fmt.Println("db end")
}
