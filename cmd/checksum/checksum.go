package main

import (
	"github.com/godcong/go-trait"
	"github.com/yinhevr/seed"
	"github.com/yinhevr/seed/model"
	"os"
	"path/filepath"
)

var log = trait.NewZapSugar()

func main() {
	args := os.Args
	dir, err := os.Getwd()
	if len(args) > 1 {
		err = nil
		dir = args[1]
	}
	if err != nil {
		log.Info("wd:", err)
		return
	}
	files := seed.GetFiles(dir)

	_ = os.MkdirAll(filepath.Join(dir, "del"), os.ModePerm)
	for _, f := range files {
		_, name := filepath.Split(f)
		sum := model.Checksum(f)
		//if name == "Yuna Sudo.jpg" {
		//	log.With("name", name, "sum", sum).Info("print name")
		//	panic(sum)
		//}
		log.With("checksum", sum).Info(f)
		if sum == "204e87534d67a17fbbaa553712ea317d845e2f27" || sum == "25b363373a120c1a5cf73f9c21fa1ca8415cfbe4" || sum == "97d573a5b0cc474eb1e95265960b4a066f3aa4b7" {
			_ = os.Rename(f, filepath.Join(dir, "del", name))
		}
	}
}
