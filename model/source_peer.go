package model

// SourcePeerDetail ...
type SourcePeerDetail struct {
	Addr string `json:"addr"`
	Peer string `json:"peer"`
}

// SourcePeer ...
type SourcePeer struct {
	//Model             `xorm:"extends"`
	*SourcePeerDetail `xorm:"extends"`
}
