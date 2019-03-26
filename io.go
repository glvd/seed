package seed

import (
	"github.com/json-iterator/go"
	"io/ioutil"
	"os"
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

func WriteJSON(path string, v interface{}) (e error) {
	bytes, e := jsoniter.Marshal(v)
	if e != nil {
		return e
	}
	e = ioutil.WriteFile(path, bytes, os.ModePerm)
	if e != nil {
		return e
	}
	return nil
}
