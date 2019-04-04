package seed

import (
	"context"
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

	video := ListVideoGet(source)
	if source.PosterPath != "" {
		s, e := rest.AddFile(source.PosterPath)
		if e != nil {
			return e
		}
		video.VideoInfo.Poster = s.Hash
	}
	log.Info(*source)
	fn := add
	if source.Slice {
		log.Debug("add slice")
		fn = addSlice
	}
	e = fn(video, source)
	if e != nil {
		return e
	}
	info := GetSourceInfo()
	log.Info(*info)

	AddSourceInfo(video, info)

	VideoListAdd(source, video)

	e = SaveVideos()
	if e != nil {
		return e
	}
	return nil
}

// GetSourceInfo ...
func GetSourceInfo() *SourceInfo {
	out, e := rest.ID()
	if e != nil {
		return &SourceInfo{}
	}
	return (*SourceInfo)(out)
}

// MustString  must string
func Mustring(val, src string) string {
	if val != "" {
		return val
	}
	return src
}

func newHLS(def *HLS) *HLS {
	if def != nil {
		def.Key = Mustring(def.Key, "")
		def.M3U8 = Mustring(def.M3U8, "media")
		def.SegmentFile = Mustring(def.SegmentFile, "media-%05d.ts")
	}

	return &HLS{
		Encrypt:     false,
		Key:         "",
		M3U8:        "media",
		SegmentFile: "media-%05d.ts",
	}
}

func addSlice(video *Video, source *VideoSource) (e error) {
	s := *source
	s.HLS = newHLS(s.HLS)
	s.Files = nil
	for _, value := range source.Files {

		e := SplitVideo(context.Background(), &s, value)
		if e != nil {
			return e
		}
	}
	e = add(video, &s)
	if e != nil {
		return e
	}

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

		group.Sliced = source.Slice
		group.HLS = source.HLS

		if dir {
			rets, e := rest.AddDir(value)
			if e != nil {
				log.Error(e)
				continue
			}
			last := len(rets) - 1
			var obj *Object
			for idx, v := range rets {
				hash = v.Hash

				if idx == last {
					obj = AddRetToLink(obj, v)
					group.Object = append(group.Object)
					continue
				}
				obj = AddRetToLinks(obj, v)
			}
			group.Object = append(group.Object, obj)

			continue
		}
		ret, e := rest.AddFile(value)
		if e != nil {
			log.Error(e)
			continue
		}
		hash = ret.Hash
		group.Object = append(group.Object, AddRetToLink(nil, ret))
	}
	group.Index = hash

	//create if null
	if video.VideoGroupList == nil {
		video.VideoGroupList = []*VideoGroup{group}
	}

	//replace if found
	for i, v := range video.VideoGroupList {
		if v.Index == group.Index {
			video.VideoGroupList[i] = group
			return nil
		}
	}
	//append if not found
	video.VideoGroupList = append(video.VideoGroupList, group)
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
