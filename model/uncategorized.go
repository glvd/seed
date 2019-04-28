package model

// Uncategorized 未分类
type Uncategorized struct {
	Name string
	Hash string
}

func init() {
	RegisterTable(&Uncategorized{})
}
