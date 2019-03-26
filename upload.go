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

	video := GetVideo(source)
	if source.PosterPath != "" {
		s, e := AddFile(source.PosterPath)
		if e != nil {
			return e
		}
		video.VideoInfo.Poster = s
	}

	fn := addNoSlick
	if source.SliceHLS {
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
	for _, value := range source.FilePath {
		s, e := AddFile(value)
		if e != nil {
			return e
		}
		ls, e := List(prefix(s))
		if e != nil {
			return e
		}

		for _, val := range ls {
			group.PlayList = append(group.PlayList, LsLinkToVideoLink(val))
		}
	}
	video.VideoGroupList = append(video.VideoGroupList, group)
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
