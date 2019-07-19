package main

import (
	"fmt"
	"github.com/yinhevr/seed"
	"log"
	"os"
	"path/filepath"
)

func main() {
	args := os.Args
	dir, err := os.Getwd()
	if len(args) > 1 {
		err = nil
		dir = args[1]
	}
	if err != nil {
		log.Println("wd:", err)
		return
	}

	files := seed.GetFiles(dir)
	for _, f := range files {
		path, name := filepath.Split(f)
		if name == "thumb.jpg" {
			continue
		}
		list := filepath.SplitList(path)
		last := len(list) - 1
		if last > 0 {
			name = list[last]
		}
		fmt.Println("from", f, "to", filepath.Join(dir, name))
		err := os.Rename(f, filepath.Join(dir, name))
		if err != nil {
			fmt.Println("error:", err)
			return
		}

	}

}
