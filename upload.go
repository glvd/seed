package seed

import (
	"golang.org/x/xerrors"
)

func prefix(s string) (ret string) {
	ret = "/ipfs/" + s
	return
}

func Upload(source *VideoSource) (e error) {
	if source == nil {
		return xerrors.New("nil source")
	}

	video := NewVideo(source)
	if source.PosterPath != "" {
		s, e := AddFile(source.PosterPath)
		if e != nil {
			return e
		}
		video.VideoInfo.Poster = s
	}

	for _, v := range source.FilePath {
		if source.Slice {
			//TODO:Slice
		}
		s, e := AddFile(v)
		if e != nil {
			return e
		}

	}

	return nil
}

func Load(path string) []*VideoSource {
	var vs []*VideoSource

	e := ReadJSON(path, &vs)
	if e != nil {
		return nil
	}
	return vs
}
