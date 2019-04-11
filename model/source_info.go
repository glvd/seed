package model

// SourceInfoDetail ...
type SourceInfoDetail struct {
	ID              string   `xorm:"-" json:"id"`
	PublicKey       string   `json:"public_key"`
	Addresses       []string `xorm:"json" json:"addresses"` //一组节点源列表
	AgentVersion    string   `json:"agent_version"`
	ProtocolVersion string   `json:"protocol_version"`
}

// SourceInfo ...
type SourceInfo struct {
	//Model             `xorm:"extends"`
	*SourceInfoDetail `xorm:"extends"`
}

// AddSourceInfo ...
func addSourceInfo(video *Video, info *SourceInfoDetail) {
	if video.SourceInfoList == nil {
		video.SourceInfoList = []*SourceInfo{{
			SourceInfoDetail: info,
		}}
		return
	}
	for idx, value := range video.SourceInfoList {
		if value.SourceInfoDetail.ID == info.ID {
			video.SourceInfoList[idx] = &SourceInfo{
				SourceInfoDetail: info,
			}
			return
		}
	}
	video.SourceInfoList = append(video.SourceInfoList, &SourceInfo{
		SourceInfoDetail: info,
	})
}
