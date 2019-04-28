package model

// Uncategorized 未分类
type Uncategorized struct {
	Name           string
	Hash           string
	IsVideo        bool
	VideoGroupList []*VideoGroup `xorm:"json" json:"video_group_list"`
}

func init() {
	RegisterTable(&Uncategorized{})
}
