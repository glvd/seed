package main

import (
	"fmt"
	"github.com/glvd/seed"
	"log"
	"os"
	"path/filepath"
	"strings"
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
	fmt.Println("load", dir)
	files := seed.GetFiles(dir)
	for _, f := range files {
		path, name := filepath.Split(f)
		if name == "thumb.jpg" {
			continue
		}
		ext := filepath.Ext(f)
		list := strings.Split(path, string(os.PathSeparator))
		fmt.Println(list)
		last := len(list) - 1
		if last > 0 {
			name = list[last]
			if name == "" {
				if last-1 >= 0 {
					name = list[last-1]
				}
			}

		}
		fmt.Println("from", f, "to", filepath.Join(dir, name))
		if name == "poster.jpg" {
			continue
		}
		err := os.Rename(f, filepath.Join(dir, name+ext))
		if err != nil {
			fmt.Println("error:", err)
			return
		}

	}

}
