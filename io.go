package seed

import (
	"github.com/json-iterator/go"
	"io/ioutil"
)

func ReadJSON(path string, v interface{}) (e error) {
	bytes, e := ioutil.ReadFile(path)
	if e != nil {
		return e
	}
	e = jsoniter.Unmarshal(bytes, v)
	if e != nil {
		return e
	}
	return nil
}
