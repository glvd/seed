package seed

import (
	"context"
	"crypto/md5"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	cmd "github.com/godcong/go-ffmpeg-cmd"

	shell "github.com/godcong/go-ipfs-restapi"
	"github.com/yinhevr/seed/model"
	"golang.org/x/xerrors"
)

func dummy(process *Process) (e error) {
	log.Info("dummy called")
	return
}

// Process ...
type Process struct {
	Seed      *Seed
	Shell     *shell.Shell
	Workspace string `json:"workspace"`
	//process   map[string]ProcessCallbackFunc
	thread  int `json:"thread"`
	ignores map[string][]byte
	//PinProc       bool   `json:"pin"`
	//Slice     bool   `json:"slice"`
	//Move      bool   `json:"move"`
	//MovePath  string `json:"move_path"`
	//JSON      bool   `json:"json"`
	//JSONPath  string `json:"json_path"`
}

func tmp(path string, name string) string {
	mp, e := filepath.Abs(path)
	if e != nil {
		mp, e = filepath.Abs(filepath.Dir(os.Args[0]))
		if e != nil {
			//ignore error
			mp, _ = os.UserHomeDir()
		}
	}
	return filepath.Join(mp, name)
}

// NewProcessSeeder ...
func NewProcessSeeder(ws string, ps ...Options) Seeder {
	process := &Process{
		Workspace: ws,
		ignores:   make(map[string][]byte, 3),
	}
	ps = append(ps, ProcessOption(process))
	return NewSeeder(ps...)
}

func prefix(s string) (ret string) {
	ret = "/ipfs/" + s
	return
}

func (p *Process) slice(unfin *model.Unfinished, format *cmd.StreamFormat, file string) (err error) {
	sa, err := cmd.FFMpegSplitToM3U8(nil, file, cmd.StreamFormatOption(format), cmd.OutputOption("tmp"))
	if err != nil {
		return err
	}
	log.Infof("%+v", sa)
	dirs, err := rest.AddDir(sa.Output)
	if err != nil {
		return err
	}

	last := unfin.SliceObject.ParseLinks(dirs)
	if last != nil {
		unfin.SliceHash = last.Hash
	}
	return nil
}

func fixPath(file string) string {
	n := strings.Replace(file, " ", "", -1)
	dir, _ := filepath.Split(n)
	//newPath := tmp("proc", name)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		log.Error(err)
	}
	err = os.Rename(file, n)
	if err != nil {
		log.Error(err)
	}
	return n
}

// Run ...
func (p *Process) Run(ctx context.Context) {
	files := p.getFiles(p.Workspace)
	log.Info(files)
	for _, oldFile := range files {
		file := fixPath(oldFile)
		log.With("old", oldFile, "new", file).Info("print filename")
		select {
		case <-ctx.Done():
			if err := ctx.Err(); err != nil {
				log.Error(err)
			}
			return
		default:
			log.With("file", file).Info("process run")
			unfin := DefaultUnfinished(file)
			object, err := rest.AddFile(file)
			if err != nil {
				log.Error(err)
				continue
			}
			//fix name and get format
			format, err := parseUnfinishedFromStreamFormat(file, unfin)
			if err != nil {
				log.Error(err)
				continue
			}
			log.Infof("%+v", format)

			unfin.Hash = object.Hash
			unfin.Object.Link = model.ObjectToVideoLink(object)
			if unfin.IsVideo {
				err := p.slice(unfin, format, file)
				if err != nil {
					log.With("split", file).Error(err)
					continue

				}
			}
			if err := model.AddOrUpdateUnfinished(unfin); err != nil {
				log.Error(err)
				continue
			}

		}
	}
	return
}

// PathMD5 ...
func PathMD5(s ...string) string {
	str := filepath.Join(s...)
	return fmt.Sprintf("%x", md5.Sum([]byte(str)))
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
		var fullPath string
		for _, name := range names {
			fullPath = filepath.Join(ws, name)
			if p.CheckIgnore(fullPath) {
				continue
			}
			tmp := p.getFiles(fullPath)
			if tmp != nil {
				files = append(files, tmp...)
			}
		}
		return files
	}
	return append(files, ws)
}

func parseUnfinishedFromStreamFormat(file string, u *model.Unfinished) (format *cmd.StreamFormat, e error) {
	format, e = cmd.FFProbeStreamFormat(file)
	if e != nil {
		return nil, e
	}

	if format.IsVideo() {
		u.IsVideo = true
		//u.Type = "video"
		u.Name = format.NameAnalyze().ToString()
		u.Sharpness = getVideoResolution(format)
	}
	return format, nil
}

// DefaultUnfinished ...
func DefaultUnfinished(name string) *model.Unfinished {
	_, file := filepath.Split(name)
	uncat := &model.Unfinished{
		Model:       model.Model{},
		Checksum:    "",
		Type:        "deprecated",
		Name:        file,
		Hash:        "",
		SliceHash:   "",
		IsVideo:     false,
		Sharpness:   "",
		Sync:        false,
		Sliced:      false,
		Encrypt:     false,
		Key:         "",
		M3U8:        "media.m3u8",
		Caption:     "",
		SegmentFile: "media-%05d.ts",
		Object:      new(model.VideoObject),
		SliceObject: new(model.VideoObject),
	}
	log.With("file", name).Info("calculate checksum")
	uncat.Checksum = model.Checksum(name)
	return uncat
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
		log.Debug("add from Unfinished")
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
	//hash := Hash(source)
	//group := parseGroup(hash, source)
	//for _, value := range source.CheckFiles {
	//	Unfinished, e := model.FindUnfinished(value, false)
	//	if e != nil {
	//		return e
	//	}
	//	group.Object = []*model.VideoObject{Unfinished.Object}
	//}

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
	//hash := Hash(source)
	//group := parseGroup(hash, source)
	//for _, value := range source.Files {
	//	info, e := os.Stat(value)
	//	if e != nil {
	//		log.Error(e)
	//		continue
	//	}
	//	dir := info.IsDir()
	//
	//	if dir {
	//		rets, e := rest.AddDir(value)
	//		if e != nil {
	//			log.Error(e)
	//			continue
	//		}
	//		last := len(rets) - 1
	//		var obj *model.VideoObject
	//		for idx, v := range rets {
	//			if idx == last {
	//				obj = model.ObjectIntoLink(obj, v)
	//				//group.Object = append(group.Object)
	//				continue
	//			}
	//			obj = model.ObjectIntoLinks(obj, v)
	//		}
	//		group.Object = append(group.Object, obj)
	//
	//		continue
	//	}
	//	ret, e := rest.AddFile(value)
	//	if e != nil {
	//		log.Error(e)
	//		continue
	//	}
	//	//hash = ret.Hash
	//	group.Object = append(group.Object, model.ObjectIntoLink(nil, ret))
	//ret
	//}

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
