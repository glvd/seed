package seed

import (
	"github.com/godcong/go-ffmpeg-cmd"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
	"os"
	"strings"
)

func prefix(s string) (ret string) {
	ret = "/ipfs/" + s
	return
}

// Upload ...
func Upload(source *VideoSource) (e error) {
	if source == nil {
		return xerrors.New("nil source")
	}

	video := GetVideo(source)
	if source.PosterPath != "" {
		s, e := rest.AddFile(source.PosterPath)
		if e != nil {
			return e
		}
		video.VideoInfo.Poster = s.Hash
	}

	fn := add
	if source.Slice {
		fn = addSlice
	}
	e = fn(video, source)
	if e != nil {
		return e
	}
	SetVideo(source, video)

	e = SaveVideos()
	if e != nil {
		return e
	}
	return nil
}

func addSlice(video *Video, source *VideoSource) (e error) {
	cmd.New()
	return nil
}

func add(video *Video, source *VideoSource) (e error) {
	group := NewVideoGroup()
	hash := ""
	for _, value := range source.Files {
		info, e := os.Stat(value)
		if e != nil {
			log.Error(e)
			continue
		}
		dir := info.IsDir()

		group.Sliced = false
		group.Sharpness = source.Sharpness
		if dir {
			rets, e := rest.AddDir(value)
			if e != nil {
				log.Error(e)
				continue
			}
			last := len(rets) - 1
			for idx, v := range rets {
				hash = v.Hash
				if idx == last {
					group.Object = AddRetToLink(group.Object, v)
				}
				group.Object = AddRetToLinks(group.Object, v)
			}

			continue
		}
		ret, e := rest.AddFile(value)
		if e != nil {
			log.Error(e)
			continue
		}

		group.Object = AddRetToLink(group.Object, ret)
		hash = ret.Hash
	}

	if video.VideoGroupList == nil {
		video.VideoGroupList = make(map[string]*VideoGroup)
	}
	video.VideoGroupList[GroupIndex(source, hash)] = group
	return nil
}

// GroupIndex ...
func GroupIndex(source *VideoSource, hash string) (s string) {
	switch strings.ToLower(source.Group) {
	case "bangumi":
		s = source.Bangumi
	case "sharpness":
		s = source.Sharpness
	case "hash":
		return hash
	default:
		s = uuid.Must(uuid.NewRandom()).String()
	}
	return
}

// Load ...
func Load(path string) []*VideoSource {
	var vs []*VideoSource
	e := ReadJSON(path, &vs)
	if e != nil {
		return nil
	}
	return vs
}
