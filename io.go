package seed

import (
	"github.com/json-iterator/go"
	"io/ioutil"
	"os"
)

// ReadJSON ...
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

// WriteJSON ...
func WriteJSON(path string, v interface{}) (e error) {
	file, e := os.OpenFile(path, os.O_WRONLY|os.O_SYNC|os.O_CREATE, os.ModePerm)
	if e != nil {
		return e
	}
	encoder := jsoniter.NewEncoder(file)
	if e = encoder.Encode(v); e != nil {
		return
	}
	return nil
}
