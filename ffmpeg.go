package seed

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	cmd "github.com/godcong/go-ffmpeg-cmd"
	"github.com/google/uuid"
	"github.com/yinhevr/seed/model"
)

const sliceM3u8FfmpegTemplate = "-y -i %s -strict -2 -c:a aac -c:v libx264 -bsf:v h264_mp4toannexb -f hls -hls_time 10 -hls_list_size 0 -hls_segment_filename %s %s"

/*
120		160 QQVGA
144		192				256
160			240 HQVGA
200				320 CGA
240		320 QVGA	360 WQVGA	384 WQVGA	400 WQVGA	432 FWQVGA (9∶5)
320			480 HVGA
360		480				640 nHD
480		640 VGA	720 WVGA or 480p	768 WVGA	800 WVGA	854 FWVGA
540						960 qHD
576						1024 WSVGA
600	750	800 SVGA			1024 WSVGA (128∶75)
640			960 DVGA	1024		1136
720		960		1152		1280 HD/WXGA
768	960	1024 XGA	1152 WXGA		1280 WXGA	1366 FWXGA
800				1280 WXGA
864		1152 XGA+	1280			1536
900		1200		1440 WXGA+		1600 HD+
960		1280 SXGA−/UVGA	1440 FWXGA+
1024	1280 SXGA		1600 WSXGA
1050		1400 SXGA+		1680 WSXGA+
1080		1440				1920 FHD
2048 DCI 2K (256∶135 ≈ 1.90∶1)	2560 UW-UXGA	3840
1152						2048 QWXGA
1200	1500	1600 UXGA		1920 WUXGA		2133		3840 (32:10)
1280			1920
1440		1920	2160 FHD+			2560 QHD/WQHD	3440 UWQHD (43∶18 = 2.38)	5120
1536		2048 QXGA
1600			2400	2560 WQXGA			3840 UW4K (12∶5 = 2.4)
1620						2880
1800				2880		3200 WQXGA+
1824			2736
1920		2560	2880	3072
2048	2560 QSXGA			3200 WQSXGA (25∶16 = 1.5625)
2160		2880	3240			3840 4K UHD
4096 DCI 4K (256∶135 ≈ 1.90∶1)	5120
2400		3200 QUXGA		3840 WQUXGA
2560			3840	4096
2880						5120 5K (UHD+)
3072		4096 HXGA
3200				5120 WHXGA
4096	5120 HSXGA			6400 WHSXGA (25∶16 = 1.5625)
4320						7680 8K UHD
4800		6400 HUXGA		7680 WHUXGA
*/
func getVideoResolution(format *cmd.StreamFormat) string {
	idx := 0
	for _, s := range format.Streams {
		if s.CodecType == "video" {
			if s.Height != nil {
				idx = getResolutionIndex(*s.Height, 0, -1)
				break
			}
		}
	}
	return strconv.FormatInt(int64(resolution[idx]), 10) + "P"

}

var resolution = []int{120, 144, 160, 200, 240, 320, 360, 480, 540, 576, 600, 640, 720, 768, 800, 864, 900, 960, 1024, 1050, 1080, 1152, 1200, 1280, 1440, 1536, 1600, 1620, 1800, 1824, 1920, 2048, 2160, 2400, 2560, 2880, 3072, 3200, 4096, 4320, 4800}

func getResolutionIndex(n int64, sta, end int) int {
	//log.Infof("%d,%d,%d", n, sta, end)
	//if int64(resolution[sta]) == n {
	//	return sta
	//}
	if end == -1 {
		end = len(resolution)
	}

	if idx := (sta + end) / 2; idx > sta {
		if int64(resolution[idx]) > n {
			return getResolutionIndex(n, sta, idx)
		}
		return getResolutionIndex(n, idx, end)
	}
	if int64(resolution[sta]) != n && sta < len(resolution)-1 {
		return sta + 1
	}
	return sta
}

// SplitVideo ...
func SplitVideo(ctx context.Context, uncat *model.Uncategorized, file string) (fp string, e error) {
	fp = filepath.Join("tmp", uuid.New().String())
	log.Debug("split path:", fp)
	fp, e = filepath.Abs(fp)
	if e != nil {
		return "", e
	}
	_ = os.MkdirAll(fp, os.ModePerm)
	format, e := cmd.FFProbeStreamFormat(file)
	if e != nil {
		return "", e
	}

	for _, s := range format.Streams {
		if s.CodecType == "video" {
			if s.Width != nil && s.Height != nil {
				getVideoResolution(format)
			}

		}
	}

	//command := cmd.NewFFMpeg()
	sf := filepath.Join(fp, uncat.SegmentFile)
	m3u8 := filepath.Join(fp, uncat.M3U8)
	args := fmt.Sprintf(sliceM3u8FfmpegTemplate, file, sf, m3u8)
	//ffmpeg -y -i $input -strict -2 -hls_segment_filename ./output/media-%03d.ts  -c:a aac -c:v copy libx264 -bsf:v h264_mp4toannexb -f hls -hls_time 10 -hls_list_size 0 ./output/m3u8

	cmd.FFMpegSpliteMedia(ctx, args)
	return "", nil
}
