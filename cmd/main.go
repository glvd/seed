package main

import (
	"flag"
	"github.com/girlvr/seed"
	"github.com/girlvr/seed/model"
	"github.com/godcong/go-trait"
	_ "github.com/mattn/go-sqlite3"

	log "github.com/sirupsen/logrus"
)

var json = flag.String("path", "seed.json", "set the path to load")
var action = flag.String("action", "process", "set action to do something")

func main() {

	defer func() {
		if e := recover(); e != nil {
			log.Panic(e)
		}
	}()

	flag.Parse()

	trait.InitRotateLog("logs/seed.log", trait.RotateLogLevel(trait.RotateLogDebug))

	switch *action {
	case "process":
		vs := seed.Load(*json)
		for _, v := range vs {
			e := seed.Upload(v)
			if e != nil {
				panic(e)
			}
		}

		log.Infof("%+v", vs[0])
	case "transfer":
		e := model.InitDB()
		if e != nil {
			panic(e)
		}
		seed.Transfer()
	}

}
