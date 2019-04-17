package seed

import (
	"context"
	"github.com/girlvr/seed/model"
	shell "github.com/godcong/go-ipfs-restapi"
	log "github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
	"os"
)

func prefix(s string) (ret string) {
	ret = "/ipfs/" + s
	return
}

// Process ...
func Process(source *VideoSource) (e error) {
	if source == nil {
		return xerrors.New("nil source")
	}
	video := &model.Video{}
	_, err := model.FindVideo(source.Bangumi, video)
	if err != nil {
		return err
	}
	parseVideoBase(video, source)
	video.SourcePeerList = nil
	video.SourceInfoList = nil
	video.VideoGroupList = nil

	log.Printf("%+v", video)
	if source.PosterPath != "" {
		object, e := rest.AddFile(source.PosterPath)
		if e != nil {
			return e
		}
		video.VideoBase.Poster = object.Hash
	}
	log.Info(*source)
	fn := add
	if source.Slice {
		log.Debug("add with slice")
		fn = addSlice
	}
	e = fn(video, source)
	if e != nil {
		return e
	}
	info := GetSourceInfo()
	log.Info(*info)

	if info.ID != "" {
		video.AddSourceInfo(info)
	}

	for _, value := range GetPeers() {
		video.AddPeers(&model.SourcePeerDetail{
			Addr: value.Addr,
			Peer: value.Peer,
		})
	}

	if err := model.AddOrUpdateVideo(video); err != nil {
		return err
	}

	//VideoListAdd(source, video)

	//if err := SaveVideos(); err != nil {
	//	return err
	//}

	return nil
}

// GetSourceInfo ...
func GetSourceInfo() *model.SourceInfoDetail {
	out, e := rest.ID()
	if e != nil {
		return &model.SourceInfoDetail{}
	}
	return (*model.SourceInfoDetail)(out)
}

// GetPeers ...
func GetPeers() []shell.SwarmConnInfo {
	swarmPeers, e := rest.SwarmPeers(context.Background())
	if e != nil {
		return nil
	}
	return swarmPeers.Peers[:50]
}

// Mustring  must string
func Mustring(val, src string) string {
	if val != "" {
		return val
	}
	return src
}

func newHLS(def *model.HLS) *model.HLS {
	if def != nil {
		def.Key = Mustring(def.Key, "")
		def.M3U8 = Mustring(def.M3U8, "media.m3u8")
		def.SegmentFile = Mustring(def.SegmentFile, "media-%05d.ts")
	}

	return &model.HLS{
		Encrypt:     false,
		Key:         "",
		M3U8:        "media.m3u8",
		SegmentFile: "media-%05d.ts",
	}
}

func addSlice(video *model.Video, source *VideoSource) (e error) {
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

func add(video *model.Video, source *VideoSource) (e error) {
	hash := Hash(source)
	group := parseGroup(hash, source)
	for _, value := range source.Files {
		info, e := os.Stat(value)
		if e != nil {
			log.Error(e)
			continue
		}
		dir := info.IsDir()

		if dir {
			rets, e := rest.AddDir(value)
			if e != nil {
				log.Error(e)
				continue
			}
			last := len(rets) - 1
			var obj *model.VideoObject
			for idx, v := range rets {
				if idx == last {
					obj = model.ObjectToLink(obj, v)
					group.Object = append(group.Object)
					continue
				}
				obj = model.ObjectToLinks(obj, v)
			}
			group.Object = append(group.Object, obj)

			continue
		}
		ret, e := rest.AddFile(value)
		if e != nil {
			log.Error(e)
			continue
		}
		//hash = ret.Hash
		group.Object = append(group.Object, model.ObjectToLink(nil, ret))
	}

	//create if null
	if video.VideoGroupList == nil {
		video.VideoGroupList = []*model.VideoGroup{group}
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

func parseGroup(index string, source *VideoSource) *model.VideoGroup {
	return &model.VideoGroup{
		Index:        index,
		Sharpness:    source.Sharpness,
		Output:       source.Output,
		Season:       source.Season,
		TotalEpisode: source.TotalEpisode,
		Episode:      source.Episode,
		Language:     source.Language,
		Caption:      source.Caption,
		Sliced:       source.Slice,
		HLS:          source.HLS,
		Object:       nil,
	}

}

// GroupIndex ...
func GroupIndex(source *VideoSource, hash string) (s string) {
	//switch strings.ToLower(source.Group) {
	//case "bangumi":
	//	s = source.Bangumi
	//case "sharpness":
	//	s = source.Sharpness
	//case "hash":
	//	return hash
	//default:
	//	s = uuid.Must(uuid.NewRandom()).String()
	//}
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
