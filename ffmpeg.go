package seed

import (
	"context"
	cmd "github.com/godcong/go-ffmpeg-cmd"
	log "github.com/sirupsen/logrus"
)

// SplitVideo ...
func SplitVideo(ctx context.Context, file string, path string) (e error) {
	command := cmd.NewFFMPEG()
	command.Split(path).Input(file).Ignore().CodecAudio(cmd.String("aac")).CodecVideo(cmd.String("libx264")).
		BitStreamFiltersVideo("h264_mp4toannexb").Format("hls").HlsTime("10").
		HlsListSize("0")
	info := make(chan string)
	cls := make(chan bool)
	go func() {
		e := command.RunContext(ctx, info, cls)
		if e != nil {
			return
		}
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
