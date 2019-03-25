package seed

import (
	"github.com/json-iterator/go"
	"io/ioutil"
)

type VideoSource struct {
	Bangumi string        `json:"bangumi"` //番号
	Path    string        `json:"path"`    //存放路径
	Poster  string        `json:"poster"`  //海报
	Role    []interface{} `json:"role"`    //主演
	Publish string        `json:"publish"` //发布
}

func ReadJSON(path string) []*VideoSource {
	var vs []*VideoSource
	bytes, e := ioutil.ReadFile(path)
	if e != nil {
		panic(e)
	}
	e = jsoniter.Unmarshal(bytes, &vs)
	if e != nil {
		panic(e)
	}
	return vs
}
