package seed

import (
	"io/ioutil"
	"path/filepath"
	"time"
)

// DaemonCallback ...
type DaemonCallback func(path string)

// DaemonStart ...
func DaemonStart(monitorPath string) {
	for {
		if lists := procList(monitorPath, func(path string) {

		}); lists != nil {
			log.Info("processing", monitorPath)

		}
		log.Info("waiting for new files")
		time.Sleep(30 * time.Second)
	}
}

func procList(monitorPath string, dc DaemonCallback) (e error) {
	infos, e := ioutil.ReadDir(monitorPath)
	if e != nil {
		panic(e)
	}
	for _, info := range infos {
		name := info.Name()
		log.Infof("name:%s ext:%s", name, filepath.Ext(name))
		if !info.IsDir() && filepath.Ext(name) != "filepart" {
			fullpath := filepath.Join(monitorPath, info.Name())
			dc(fullpath)
		}
		if name == "success" {
			continue
		}
	}
	return nil
}
