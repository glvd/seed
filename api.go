package seed

import (
	"github.com/ipfs/go-ipfs-api"
	log "github.com/sirupsen/logrus"
)

var sh = shell.NewShell("localhost:5001")

func AddDir(dir string) (e error) {
	s, e := sh.AddDir(dir)
	if e != nil {
		return e
	}
	log.Info(s)
	return nil
}

func Dir(path string) (e error) {
	links, e := sh.List(path)
	if e != nil {
		return e
	}
	log.Infof("%+v", links)
	return nil
}
