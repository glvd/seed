package model

// SourcePeerDetail ...
type SourcePeerDetail struct {
	Addr string `json:"addr"`
	Peer string `json:"peer"`
}

// SourcePeer ...
type SourcePeer struct {
	Model            `xorm:"extends"`
	Index            int
	SourcePeerDetail `xorm:"extends"`
}

func init() {
	RegisterTable(SourcePeer{})
}
