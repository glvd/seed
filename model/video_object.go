package model

// Link ...
type VideoLink struct {
	Hash string `json:"hash"`
	Name string `json:"name"`
	Size uint64 `json:"size"`
	Type int    `json:"type"`
}

// Object ...
type VideoObject struct {
	Model     `xorm:"extends"`
	Links     []*VideoLink `json:"links,omitempty"`
	VideoLink `json:",inline"`
}
