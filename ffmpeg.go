package seed

import (
	"context"
	cmd "github.com/godcong/go-ffmpeg-cmd"
	"github.com/google/uuid"
	"github.com/yinhevr/seed/model"
	"os"
	"path/filepath"
)

// SplitVideo ...
func SplitVideo(ctx context.Context, hls *model.HLS, file string) (fp string, e error) {
	fp = filepath.Join("tmp", uuid.New().String())
	log.Debug("split path:", fp)
	fp, e = filepath.Abs(fp)
	if e != nil {
		return "", e
	}
	_ = os.MkdirAll(fp, os.ModePerm)

	command := cmd.NewFFMPEG()
	command.Split(fp).Strict().
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
	}()
	for {
		select {
		case v := <-info:
			if v != "" {
				log.Info(v)
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
