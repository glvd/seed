package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	cmd "github.com/godcong/go-ffmpeg-cmd"
	"github.com/yinhevr/seed"
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
	ctx := cmd.FFmpegContext()
	for _, f := range files {
		ext := filepath.Ext(f)
		if strings.ToUpper(ext) == ".ISO" || strings.ToUpper(ext) == ".WMV" {
			log.Println("transfer:", f)
			input := fmt.Sprintf("-y -i %s %s", f, f+".mp4")
			err := cmd.FFMpegRun(ctx, input)
			if err != nil {
				log.Println(err)
				return
			}
		}
	}

}
