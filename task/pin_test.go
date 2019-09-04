package task_test

import (
	"github.com/glvd/seed"
	"github.com/glvd/seed/model"
	"log"
	"testing"
)

// TestPin ...
func TestPin(t *testing.T) {
	db, e := model.InitSQLite3("data.db")
	if e != nil {
		log.Fatal(e)
	}
	sdb := seed.NewDatabase(db, seed.DatabaseShowSQLArg())

	seed := seed.NewSeed(sdb)

	seed.Register()

	//seed.Workspace = "D:\\videoall"
	//seed.AfterInit()
	seed.Start()
	seed.Wait()
}
