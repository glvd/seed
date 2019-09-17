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

func TestPinSync_Call_Unfinished(t *testing.T) {
	seeder := seed.NewSeed()

	engine, e := model.InitSQLite3("0916.db")
	if e != nil {
		t.Error(e)
	}
	database := seed.NewDatabase(engine)
	database.RegisterSync(model.Video{}, model.Pin{}, model.Unfinished{})
	api := seed.NewAPI("/ip4/127.0.0.1/tcp/5001")
	seeder.Register(database, api)
	pin := task.NewPin()
	pin.Type = task.PinTypeSync
	pin.Table = task.PinTableVideo
	pin.From = "/ip4/192.168.1.13/tcp/14001/ipfs/QmXNZRTd54Zvarf4sswVvUUnpb4gPQNAhFViozVgG8uwri"
	skip := []string{"poster", "thumb", "video"}
	for _, s := range skip {
		pin.SkipType = append(pin.SkipType, s)
	}

	seeder.Start()
	seeder.AddTasker(pin)
	seeder.Wait()

}

func TestPinSync_Call_Video(t *testing.T) {
	seeder := seed.NewSeed()

	engine, e := model.InitSQLite3("0916.db")
	if e != nil {
		t.Error(e)
	}
	database := seed.NewDatabase(engine, seed.DatabaseShowSQLArg())
	database.RegisterSync(model.Video{}, model.Pin{}, model.Unfinished{})
	api := seed.NewAPI("/ip4/127.0.0.1/tcp/5001")
	seeder.Register(database, api)
	pin := task.NewPin()
	pin.Type = task.PinTypeSync
	pin.Table = task.PinTableVideo
	pin.From = "/ip4/192.168.1.13/tcp/14001/ipfs/QmXNZRTd54Zvarf4sswVvUUnpb4gPQNAhFViozVgG8uwri"
	skip := []string{"thumb", "video", "slice"}
	for _, s := range skip {
		pin.SkipType = append(pin.SkipType, s)
	}

	seeder.Start()
	seeder.AddTasker(pin)
	seeder.Wait()

}
func TestPinSync_Call(t *testing.T) {
	seeder := seed.NewSeed()

	engine, e := model.InitSQLite3("0916.db")
	if e != nil {
		t.Error(e)
	}
	database := seed.NewDatabase(engine)
	database.RegisterSync(model.Video{}, model.Pin{}, model.Unfinished{})

	api := seed.NewAPI("/ip4/127.0.0.1/tcp/5001")
	seeder.Register(database, api)
	pin := task.NewPin()
	pin.Type = task.PinTypeSync
	pin.Table = task.PinTablePin
	pin.From = "/ip4/192.168.1.13/tcp/14001/ipfs/QmXNZRTd54Zvarf4sswVvUUnpb4gPQNAhFViozVgG8uwri"
	skip := []string{""}
	for _, s := range skip {
		pin.SkipType = append(pin.SkipType, s)
	}

	seeder.Start()
	seeder.AddTasker(pin)
	seeder.Wait()
}
