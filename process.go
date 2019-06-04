package seed

import (
	"context"
	"crypto/md5"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	shell "github.com/godcong/go-ipfs-restapi"
	"github.com/yinhevr/seed/model"
	"golang.org/x/xerrors"
	"gopkg.in/urfave/cli.v2"
)

// CallbackFunc ...
type CallbackFunc func(interface{}) error

// Process ...
type Process struct {
	Workspace string `json:"workspace"`
	Pin       bool   `json:"pin"`
	Slice     bool   `json:"slice"`
	Move      bool   `json:"move"`
	MovePath  string `json:"move_path"`
	JSON      bool   `json:"json"`
	JSONPath  string `json:"json_path"`
	Before    CallbackFunc
	After     CallbackFunc
	ignores   map[string][]byte
}

func initIgnore() map[string][]byte {
	return make(map[string][]byte, 3)
}

// NewProcess ...
func NewProcess() *Process {
	return &Process{
		Pin:      false,
		Slice:    true,
		Move:     true,
		MovePath: "tmp",
		JSON:     false,
		JSONPath: "seed.json",
		ignores:  initIgnore(),
	}
}

func prefix(s string) (ret string) {
	ret = "/ipfs/" + s
	return
}

func isVideo(filename string) bool {
	vlist := []string{
		"mkv", ".mp4", ".mpg", ".mpeg", ".avi", ".rm", ".rmvb", ".mov", ".wmv", ".asf", ".dat", ".asx", ".wvx", ".mpe", ".mpa",
	}
	for _, v := range vlist {
		if path.Ext(filename) == v {
			return true
		}
	}
	return false
}

// CmdProcess ...
func CmdProcess(app *cli.App) *cli.Command {
	flags := append(app.Flags,
		&cli.BoolFlag{
			Name:        "quick",
			Aliases:     []string{"q"},
			Usage:       "process data to ipfs skip read json",
			Value:       false,
			Destination: nil,
		})

	return &cli.Command{
		Name:    "process",
		Aliases: []string{"P"},
		Usage:   "",
		Action: func(context *cli.Context) error {
			log.Info("process call")

			if context.Bool("q") {
				//QuickProcess()
			}
			//ProcessVideo()

			return nil
		},
		Subcommands: nil,
		Flags:       flags,
	}
}

// Run ...
func (p *Process) Run() (err error) {
	p.Ignore()
	if p.Before != nil {
		if err = p.Before(p); err != nil {
			return err
		}
	}
	p.getFiles(p.Workspace)
	//uncat := model.Uncategorized{}

	return nil
}

// PathMD5 ...
func PathMD5(s ...string) string {
	str := filepath.Join(s...)
	return fmt.Sprintf("%x", md5.Sum([]byte(str)))
}

// Ignore ...
func (p *Process) Ignore() {
	p.ignores[PathMD5(p.Workspace, p.MovePath)] = nil
}

// CheckIgnore ...
func (p *Process) CheckIgnore(name string) (b bool) {
	if p.ignores == nil {
		return false
	}
	_, b = p.ignores[PathMD5(name)]
	return
}

func (p *Process) getFiles(ws string) (files []string) {
	info, e := os.Stat(ws)
	if e != nil {
		return nil
	}
	if info.IsDir() {
		file, e := os.Open(ws)
		if e != nil {
			return nil
		}
		defer file.Close()
		names, e := file.Readdirnames(-1)
		if e != nil {
			return nil
		}
		var fullpath string
		for _, name := range names {
			fullpath = filepath.Join(ws, name)
			if p.CheckIgnore(fullpath) {
				continue
			}
			tmp := p.getFiles(fullpath)
			if tmp != nil {
				files = append(files, tmp...)
			}
		}
		return files
	}
	files = append(files, ws)
	return
}

// QuickProcess ...
func QuickProcess(pathname string, needPin bool) (e error) {
	info, e := os.Stat(pathname)
	if e != nil {
		return e
	}
	b := info.IsDir()
	if b {
		file, e := os.Open(pathname)
		if e != nil {
			return e
		}
		defer file.Close()
		names, e := file.Readdirnames(-1)
		if e != nil {
			return e
		}
		for _, value := range names {
			uncat := model.Uncategorized{
				Name:    value,
				Type:    "other",
				Hash:    "",
				IsVideo: false,
				Object:  nil,
			}
			file := filepath.Join(pathname, value)
			fileinfo, e := os.Stat(file)
			if e != nil {
				log.Error(e)
				continue
			}
			if fileinfo.IsDir() {
				log.Error(value, " continue with dir")
				continue
			}
			uncat.Checksum = model.Checksum(file)
			log.Infof("add [%s]:%s", uncat.Checksum, file)
			object, e := rest.AddFile(file)
			if e != nil {
				log.Errorf("add file error:%+v", object)
				continue
			}
			log.Info("added", uncat.Checksum)
			uncat.Hash = object.Hash
			uncat.Object = append(uncat.Object, model.ObjectToLink(nil, object))
			uncat.IsVideo = isVideo(value)
			if uncat.IsVideo {
				uncat.Type = "video"
			}
			e = model.AddOrUpdateUncategorized(&uncat)
			if e != nil {
				log.Errorf("insert uncategorized error:%+v", e)
				continue
			}
			log.Info("inserted table", uncat)
			if uncat.IsVideo {
				uncatvideo := model.Uncategorized{
					Model:    model.Model{},
					Checksum: uncat.Checksum,
					Type:     "m3u8",
					Name:     value,
					Hash:     "",
					IsVideo:  uncat.IsVideo,
					Object:   nil,
				}
				file, e := SplitVideo(context.Background(), nil, file)
				if e != nil {
					log.Errorf("split file error:%+v", object)
					continue
				}
				log.Info(file)
				rets, e := rest.AddDir(file)
				if e != nil {
					log.Errorf("ipfs add file error:%+v", object)
					continue
				}
				last := len(rets) - 1
				var obj *model.VideoObject
				for idx, v := range rets {
					if idx == last {
						obj = model.ObjectToLink(obj, v)
						uncatvideo.Hash = obj.Link.Hash
						continue
					}
					obj = model.ObjectToLinks(obj, v)
				}
				log.Infof("%+v", *obj)
				uncatvideo.Object = append(uncatvideo.Object, obj)
				log.Info("inserted video table", uncatvideo)
				if needPin {
					pin(nil, uncatvideo.Hash, func(hash string) {
						uncatvideo.Sync = true
						e = model.AddOrUpdateUncategorized(&uncatvideo)
						if e != nil {
							log.Errorf("insert uncategorized error:%+v", e)
						}
					})
					continue
				}
				e = model.AddOrUpdateUncategorized(&uncatvideo)
				if e != nil {
					log.Errorf("insert uncategorized error:%+v", e)
					continue
				}
			}
			log.Info("move file", file)
			if err := moveSuccess(file); err != nil {
				return err
			}

		}

	}
	return nil
}

func moveSuccess(file string) (e error) {
	dir, name := filepath.Split(file)
	newPath := filepath.Join(dir, "success")
	_ = os.MkdirAll(newPath, os.ModePerm)
	newPathFile := filepath.Join(newPath, name)
	return os.Rename(file, newPathFile)
}

// ProcessVideo ...
func ProcessVideo(source *VideoSource) (e error) {
	if source == nil {
		return xerrors.New("nil source")
	}
	if source.Bangumi == "" {
		return nil
	}
	source.Bangumi = strings.ToUpper(source.Bangumi)
	video := &model.Video{}
	_, err := model.FindVideo(source.Bangumi, video, false)
	if err != nil {
		return err
	}
	parseVideoBase(video, source)

	log.Infof("%+v", video)

	video.PosterHash = addPosterHash(source)
	log.Info(*source)

	fn := add
	if source.CheckFiles != nil {
		log.Debug("add from uncategorized")
		fn = addChecksum
	} else if source.Slice {
		log.Debug("add with slice")
		fn = addSlice
	}
	e = fn(video, source)
	if e != nil {
		return e
	}
	info := GetSourceInfo()
	log.Info(*info)

	//if info.ID != "" {
	//	video.AddSourceInfo(info)
	//}
	//
	//for _, value := range GetPeers() {
	//	video.AddPeers(&model.SourcePeerDetail{
	//		Addr: value.Addr,
	//		Peer: value.Peer,
	//	})
	//}
	return model.AddOrUpdateVideo(video)
}

func addPosterHash(source *VideoSource) string {
	if source.PosterPath != "" {
		object, e := rest.AddFile(source.PosterPath)
		if e != nil {
			log.Error("add poster error:", e)
			return ""
		}
		return object.Hash
	}
	return source.Poster
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
	size := len(swarmPeers.Peers)
	if size == 0 {
		return nil
	}
	if size > 50 {
		size = 50
	}

	return swarmPeers.Peers[:size]
}

// MustString  must string
func MustString(val, src string) string {
	if val != "" {
		return val
	}
	return src
}

// hls ...
func initHLS(source *VideoSource) {
	if source != nil {
		source.Key = MustString(source.Key, "")
		source.M3U8 = MustString(source.M3U8, "media.m3u8")
		source.SegmentFile = MustString(source.SegmentFile, "media-%05d.ts")
	}
}

func addSlice(video *model.Video, source *VideoSource) (e error) {
	initHLS(source)
	source.Files = nil
	for _, value := range source.Files {
		file, e := SplitVideo(context.Background(), nil, value)
		if e != nil {
			return e
		}
		source.Files = append(source.Files, file)
	}
	e = add(video, source)
	if e != nil {
		return e
	}

	return nil
}

func addChecksum(video *model.Video, source *VideoSource) (e error) {
	hash := Hash(source)
	group := parseGroup(hash, source)
	for _, value := range source.CheckFiles {
		uncategorized, e := model.FindUncategorized(value, false)
		if e != nil {
			return e
		}
		group.Object = uncategorized.Object
	}

	//create if null
	//if video.VideoGroupList == nil {
	//	video.VideoGroupList = []*model.VideoGroup{group}
	//}
	//
	////replace if found
	//for i, v := range video.VideoGroupList {
	//	if v.Index == group.Index {
	//		video.VideoGroupList[i] = group
	//		return nil
	//	}
	//}
	////append if not found
	//video.VideoGroupList = append(video.VideoGroupList, group)
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
					//group.Object = append(group.Object)
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
	//if video.VideoGroupList == nil {
	//	video.VideoGroupList = []*model.VideoGroup{group}
	//}
	//
	////replace if found
	//for i, v := range video.VideoGroupList {
	//	if v.Index == group.Index {
	//		video.VideoGroupList[i] = group
	//		return nil
	//	}
	//}
	////append if not found
	//video.VideoGroupList = append(video.VideoGroupList, group)
	return nil
}

func parseGroup(index string, source *VideoSource) *model.VideoGroup {
	return &model.VideoGroup{
		Index: index,
		//Sharpness:    source.Sharpness,
		//Output:       source.Output,
		//Season:       source.Season,
		//TotalEpisode: source.TotalEpisode,
		//Episode:      source.Episode,
		//Language:     source.Language,
		//Caption:      source.Caption,
		Sliced: source.Slice,
		//HLS:          *hls(source.HLS),
		Object: nil,
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
