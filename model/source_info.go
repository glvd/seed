package model

// SourceInfoDetail ...
type SourceInfoDetail struct {
	ID              string   `xorm:"source_id" json:"id"`
	PublicKey       string   `json:"public_key"`
	Addresses       []string `xorm:"json" json:"addresses"` //一组节点源列表
	AgentVersion    string   `json:"agent_version"`
	ProtocolVersion string   `json:"protocol_version"`
}

// SourceInfo ...
type SourceInfo struct {
	Model      `xorm:"extends"`
	SourceInfo SourceInfoDetail `xorm:"extends"`
}

func init() {
	RegisterTable(SourceInfo{})
}
