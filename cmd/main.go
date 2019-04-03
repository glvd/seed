package main

import (
	"flag"
	"github.com/girlvr/seed"
	"github.com/godcong/go-trait"
	log "github.com/sirupsen/logrus"
)

var json = flag.String("path", "seed.json", "set the path to load")

func main() {

	defer func() {
		if e := recover(); e != nil {
			log.Panic(e)
		}
	}()

	flag.Parse()

	trait.InitRotateLog("logs/seed.log", trait.RotateLogLevel(trait.RotateLogDebug))

	vs := seed.Load(*json)
	for _, v := range vs {
		e := seed.Upload(v)
		if e != nil {
			log.Error(e)
		}
	}

	log.Infof("%+v", vs[0])
}
