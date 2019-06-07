package model

import shell "github.com/godcong/go-ipfs-restapi"

// VideoLink ...
type VideoLink struct {
	Hash string `json:"hash"`
	Name string `json:"name"`
	Size uint64 `json:"size"`
	Type int    `json:"type"`
}

// VideoObject ...
type VideoObject struct {
	Links []*VideoLink `json:"links,omitempty"`
	Link  *VideoLink   `xorm:"extends"  json:",inline"`
}

// ObjectToVideoLink ...
func ObjectToVideoLink(obj *shell.Object) *VideoLink {
	return &VideoLink{
		Hash: obj.Hash,
		Name: obj.Name,
		Size: obj.Size,
		Type: -1,
	}
}

// ParseLinks ...
func (obj *VideoObject) ParseLinks(links []*shell.Object) *VideoLink {
	last := len(links) - 1
	for i, link := range links {
		if i == last || last == 0 {
			obj.Link = ObjectToVideoLink(link)
			break
		}
		obj.Links = append(obj.Links, ObjectToVideoLink(link))
	}
	return obj.Link
}

//
//// ObjectIntoLinks ...
//func ObjectIntoLinks(obj *VideoObject, ret *shell.Object) *VideoObject {
//	if obj != nil {
//		obj.Links = append(obj.Links, &VideoLink{
//			Hash: ret.Hash,
//			Name: ret.Name,
//			Size: ret.Size,
//			Type: 2,
//		})
//		return obj
//	}
//	return &VideoObject{
//		Links: []*VideoLink{
//			{
//				Hash: ret.Hash,
//				Name: ret.Name,
//				Size: ret.Size,
//				Type: 2,
//			},
//		},
//	}
//}
