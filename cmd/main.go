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

	trait.InitRotateLog("logs/seed.log")

	vs := seed.ReadJSON(*json)
	for _, v := range vs {
		s, e := seed.AddDir(v.Path)
		if e != nil {
			//continue next
			log.Error(e)
		}
		log.Info("dir:", s)
	}
	log.Infof("%+v", vs)
}
