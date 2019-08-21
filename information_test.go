package seed_test

import (
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

	s.Wait()
}

// TestInformation2 ...
func TestInformation2(t *testing.T) {
	info := seed.NewInformation()
	info.ResourcePath = "D:\\videoall\\images"
	info.Path = "D:\\videoall\\video4.json"
	info.InfoType = seed.InfoTypeBSON
	sdb := seed.NewDatabase(model.MustDatabase(model.InitSQLite3("test.db")))
	sdb.RegisterSync(model.Video{}, model.Pin{}, model.Unfinished{})
	s := seed.NewSeed(info, sdb)

	s.Start()

	s.Wait()
}
