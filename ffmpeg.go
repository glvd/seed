package seed

import (
	"context"
	cmd "github.com/godcong/go-ffmpeg-cmd"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/yinhevr/seed/model"
	"os"
	"path/filepath"
)

// SplitVideo ...
func SplitVideo(ctx context.Context, hls *model.HLS, file string) (files []string, e error) {
	path := filepath.Join("tmp", uuid.New().String())
	log.Debug("split path:", path)
	path, e = filepath.Abs(path)
	if e != nil {
		return nil, e
	}
	_ = os.MkdirAll(path, os.ModePerm)

	command := cmd.NewFFMPEG()
	command.Split(path).Strict().
		HlsSegmentFilename(hls.SegmentFile).Output(hls.M3U8).
		Input(file).Ignore().
		CodecAudio(cmd.String("aac")).CodecVideo(cmd.String("libx264")).
		BitStreamFiltersVideo("h264_mp4toannexb").Format("hls").HlsTime("10").
		HlsListSize("0")
	info := make(chan string)
	cls := make(chan bool)
	go func() {
		e := command.RunContext(ctx, info, cls)
		if e != nil {
			return
		}
		files = append(files, path)
	}()
	for {
		select {
		case v := <-info:
			if v != "" {
				log.Print(v)
			}
		case c := <-cls:
			if c == true {
				close(info)
				return
			}
		case <-ctx.Done():
			log.Debug("done")
			log.Debug(ctx.Err())
			return
		default:
			//log.Println("waiting:...")
		}
	}
}
