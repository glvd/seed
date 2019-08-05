package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"regexp"

	"github.com/glvd/seed"
)

func main() {
	args := os.Args
	dir, err := os.Getwd()
	if len(args) > 1 {
		err = nil
		dir = args[1]
	}
	if err != nil {
		log.Println("wd:", err)
		return
	}
	file, err := os.Open(dir)
	if err != nil {
		log.Println(err)
		return
	}
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return
	}
	fixed := fixBson(bytes)
	openFile, err := os.OpenFile("fixed.json", os.O_CREATE|os.O_SYNC|os.O_RDWR, os.ModePerm)
	if err != nil {
		log.Println(err)
		return
	}
	_, e := openFile.Write(fixed)
	if e != nil {
		log.Println(e)
		return
	}

	var vs []*seed.VideoSource
	err = json.Unmarshal(fixed, &vs)
	if err != nil {
		log.Println(err)
		return
	}
	i := len(vs)
	log.Println("printed:", i)

}

func fixBson(s []byte) []byte {
	reg := regexp.MustCompile(`("_id")[ ]*[:][ ]*(ObjectId\(")[\w]{24}("\))[ ]*(,)[ ]*`)
	return reg.ReplaceAll(s, []byte(" "))
}
