package seed

import (
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

	fn := addNoSlick
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
	//TODO: slice
	return nil
}

func addNoSlick(video *Video, source *VideoSource) (e error) {
	group := NewVideoGroup()
	for _, value := range source.Files {
		info, e := os.Stat(value)
		if e != nil {
			log.Error(e)
			continue
		}
		dir := info.IsDir()
		if dir {
			s, e := rest.AddDirList(value)
			if e != nil {
				log.Error(e)
				continue
			}
			group.Object = s
			continue
		}
		ret, e := rest.AddFile(value)
		if e != nil {
			return e
		}
		group.Object = AddRetToObject(ret)
	}

	if video.VideoGroupList == nil {
		video.VideoGroupList = make(map[string]*VideoGroup)
	}
	video.VideoGroupList[GroupIndex(source)] = group
	return nil
}

func GroupIndex(source *VideoSource) (s string) {
	switch strings.ToLower(source.Group) {
	case "bangumi":
		s = source.Bangumi
	case "sharpness":
		s = source.Sharpness
	default:
		s = uuid.Must(uuid.NewRandom()).String()
	}
	return
}

func Load(path string) []*VideoSource {
	var vs []*VideoSource
	e := ReadJSON(path, &vs)
	if e != nil {
		return nil
	}
	return vs
}
