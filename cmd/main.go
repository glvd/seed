package main

import (
	log "github.com/sirupsen/logrus"
)

func main() {
	defer func() {
		if e := recover(); e != nil {
			log.Panic(e)
		}
	}()
}
