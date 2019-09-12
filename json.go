package seed

import (
	"encoding/json"
	"os"
)

func jsonWrite(path string, v interface{}) error {
	f, e := os.OpenFile(path, os.O_CREATE|os.O_SYNC|os.O_RDWR|os.O_TRUNC, os.ModePerm)
	if e != nil {
		return e
	}
	encoder := json.NewEncoder(f)
	return encoder.Encode(v)
}

func jsonRead(path string, v interface{}) error {

	f, e := os.Open(path)
	if e != nil {
		return e
	}
	decoder := json.NewDecoder(f)
	return decoder.Decode(v)
}
