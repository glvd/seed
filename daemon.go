package seed

import (
	"io/ioutil"
	"path/filepath"
	"time"
)

// DaemonStart ...
func DaemonStart(monitorPath string) {
	for {
		if lists := getList(monitorPath); lists != nil {
			log.Info("processing", monitorPath)
			for _, list := range lists {
				moveSuccess(list)
			}
		}
		log.Info("waiting for new files")
		time.Sleep(30 * time.Second)
	}

}

func getList(monitorPath string) (files []string) {
	infos, e := ioutil.ReadDir(monitorPath)
	if e != nil {
		panic(e)
	}
	for _, info := range infos {
		if !info.IsDir() {
			name := info.Name()
			log.Info("name:", name)
			log.Info("ext:", filepath.Ext(name))
			fullpath := filepath.Join(monitorPath, info.Name())
			log.Info("fullpath:", fullpath)

		}
	}
	return files
}
