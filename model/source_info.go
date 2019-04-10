package model

// SourceInfo ...
type SourceInfo struct {
	ID              string   `json:"id"`
	PublicKey       string   `json:"public_key"`
	Addresses       []string `xorm:"json" json:"addresses"` //一组节点源列表
	AgentVersion    string   `json:"agent_version"`
	ProtocolVersion string   `json:"protocol_version"`
}
