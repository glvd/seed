package main

import (
	"github.com/girlvr/seed"
	log "github.com/sirupsen/logrus"
)

func main() {
	defer func() {
		if e := recover(); e != nil {
			log.Panic(e)
		}
	}()

	vs := seed.ReadJSON("D:\\workspace\\goproject\\seed")
	log.Infof("%+v", vs)
}
