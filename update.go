package seed

import (
	"github.com/yinhevr/seed/model"
	"golang.org/x/xerrors"
)

// Update ...
func Update(index string, source *VideoSource) (e error) {
	if source == nil {
		return xerrors.New("nil source")
	}
	video := &model.Video{}
	b, err := model.FindVideo(source.Bangumi, video, false)
	if err != nil || !b {
		return xerrors.New("nil video")
	}
	hash := Hash(source)
	if hash == index {
		//success with no changes
		return nil
	}

	group := parseGroup(hash, source)
	for idx := range video.VideoGroupList {
		if video.VideoGroupList[idx].Index == index {
			objBak := video.VideoGroupList[idx].Object
			video.VideoGroupList[idx] = group
			video.VideoGroupList[idx].Object = objBak
			break
		}
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

	return model.AddOrUpdateVideo(video)
}
