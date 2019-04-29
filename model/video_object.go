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
	Link  VideoLink    `json:",inline"`
}

// LinkObjectToObject ...
func LinkObjectToObject(obj interface{}) *VideoObject {
	if v, b := obj.(*VideoObject); b {
		return v
	}
	return &VideoObject{}
}

// ObjectToLink ...
func ObjectToLink(obj *VideoObject, ret *shell.Object) *VideoObject {
	if obj != nil {
		obj.Link.Hash = ret.Hash
		obj.Link.Name = ret.Name
		obj.Link.Size = ret.Size
		obj.Link.Type = 2
		return obj
	}
	return &VideoObject{
		Link: VideoLink{
			Hash: ret.Hash,
			Name: ret.Name,
			Size: ret.Size,
			Type: 2,
		},
	}
}

// ObjectToLinks ...
func ObjectToLinks(obj *VideoObject, ret *shell.Object) *VideoObject {
	if obj != nil {
		obj.Links = append(obj.Links, &VideoLink{
			Hash: ret.Hash,
			Name: ret.Name,
			Size: ret.Size,
			Type: 2,
		})
		return obj
	}
	return &VideoObject{
		Links: []*VideoLink{
			{
				Hash: ret.Hash,
				Name: ret.Name,
				Size: ret.Size,
				Type: 2,
			},
		},
	}
}
