package seed

import (
	"context"
	"fmt"
	cmd "github.com/godcong/go-ffmpeg-cmd"
	"github.com/google/uuid"
	"github.com/yinhevr/seed/model"
	"os"
	"path/filepath"
)

// StreamFormat ...
type StreamFormat struct {
	Streams []Stream `json:"streams"`
	Format  Format   `json:"format"`
}

// Format ...
type Format struct {
	Filename       string     `json:"filename"`
	NbStreams      int64      `json:"nb_streams"`
	NbPrograms     int64      `json:"nb_programs"`
	FormatName     string     `json:"format_name"`
	FormatLongName string     `json:"format_long_name"`
	StartTime      string     `json:"start_time"`
	Duration       string     `json:"duration"`
	Size           string     `json:"size"`
	BitRate        string     `json:"bit_rate"`
	ProbeScore     int64      `json:"probe_score"`
	Tags           FormatTags `json:"tags"`
}

// FormatTags ...
type FormatTags struct {
	MajorBrand       string `json:"major_brand"`
	MinorVersion     string `json:"minor_version"`
	CompatibleBrands string `json:"compatible_brands"`
	Encoder          string `json:"encoder"`
}

// Stream ...
type Stream struct {
	Index              int64            `json:"index"`
	CodecName          string           `json:"codec_name"`
	CodecLongName      string           `json:"codec_long_name"`
	Profile            string           `json:"profile"`
	CodecType          string           `json:"codec_type"`
	CodecTimeBase      string           `json:"codec_time_base"`
	CodecTagString     string           `json:"codec_tag_string"`
	CodecTag           string           `json:"codec_tag"`
	Width              *int64           `json:"width,omitempty"`
	Height             *int64           `json:"height,omitempty"`
	CodedWidth         *int64           `json:"coded_width,omitempty"`
	CodedHeight        *int64           `json:"coded_height,omitempty"`
	HasBFrames         *int64           `json:"has_b_frames,omitempty"`
	SampleAspectRatio  *string          `json:"sample_aspect_ratio,omitempty"`
	DisplayAspectRatio *string          `json:"display_aspect_ratio,omitempty"`
	PixFmt             *string          `json:"pix_fmt,omitempty"`
	Level              *int64           `json:"level,omitempty"`
	ColorRange         *string          `json:"color_range,omitempty"`
	ColorSpace         *string          `json:"color_space,omitempty"`
	ColorTransfer      *string          `json:"color_transfer,omitempty"`
	ColorPrimaries     *string          `json:"color_primaries,omitempty"`
	ChromaLocation     *string          `json:"chroma_location,omitempty"`
	Refs               *int64           `json:"refs,omitempty"`
	IsAVC              *string          `json:"is_avc,omitempty"`
	NalLengthSize      *string          `json:"nal_length_size,omitempty"`
	RFrameRate         string           `json:"r_frame_rate"`
	AvgFrameRate       string           `json:"avg_frame_rate"`
	TimeBase           string           `json:"time_base"`
	StartPts           int64            `json:"start_pts"`
	StartTime          string           `json:"start_time"`
	DurationTs         int64            `json:"duration_ts"`
	Duration           string           `json:"duration"`
	BitRate            string           `json:"bit_rate"`
	BitsPerRawSample   *string          `json:"bits_per_raw_sample,omitempty"`
	NbFrames           string           `json:"nb_frames"`
	Disposition        map[string]int64 `json:"disposition"`
	Tags               StreamTags       `json:"tags"`
	SampleFmt          *string          `json:"sample_fmt,omitempty"`
	SampleRate         *string          `json:"sample_rate,omitempty"`
	Channels           *int64           `json:"channels,omitempty"`
	ChannelLayout      *string          `json:"channel_layout,omitempty"`
	BitsPerSample      *int64           `json:"bits_per_sample,omitempty"`
	MaxBitRate         *string          `json:"max_bit_rate,omitempty"`
}

// StreamTags ...
type StreamTags struct {
	Language    string `json:"language"`
	HandlerName string `json:"handler_name"`
}

// SplitVideo ...
func SplitVideo(ctx context.Context, hls *model.HLS, file string) (fp string, e error) {
	fp = filepath.Join(os.TempDir(), uuid.New().String())
	log.Debug("split path:", fp)
	fp, e = filepath.Abs(fp)
	if e != nil {
		return "", e
	}
	_ = os.MkdirAll(fp, os.ModePerm)

	command := cmd.NewFFMpeg()

	sf := filepath.Join(fp, hls.SegmentFile)
	m3u8 := filepath.Join(fp, hls.M3U8)
	args := fmt.Sprintf("ffmpeg -y -i %s -strict -2 -c:a aac -c:v copy libx264 -bsf:v h264_mp4toannexb -f hls -hls_time 10 -hls_list_size 0 -hls_segment_filename %s %s", file, sf, m3u8)
	//ffmpeg -y -i $input -strict -2 -hls_segment_filename ./output/media-%03d.ts  -c:a aac -c:v copy libx264 -bsf:v h264_mp4toannexb -f hls -hls_time 10 -hls_list_size 0 ./output/m3u8
	command.SetArgs(args)
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
