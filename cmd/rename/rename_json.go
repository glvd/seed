package main

import (
	"encoding/json"
	"github.com/godcong/go-trait"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

var log2 = trait.NewZapSugar()

func main() {
	args := os.Args
	dir, err := os.Getwd()
	if len(args) > 1 {
		err = nil
		dir = args[1]
	}
	if err != nil {
		log2.Info("wd:", err)
		return
	}

	bytes, err := ioutil.ReadFile(dir)
	if err != nil {
		panic(err)
	}

	var role Role
	err = json.Unmarshal(fixBson(bytes), &role)
	if err != nil {
		panic(err)
	}
	path, _ := filepath.Split(dir)
	for _, ele := range ([]RoleElement)(role) {
		ext := filepath.Ext(ele.Avatar)
		log2.With("from", ele.Avatar, "to", ele.Name+ext).Info("rename")
		old := filepath.Join(path, ele.Avatar)
		new := filepath.Join(path, ele.Name+ext)
		_ = os.Rename(old, new)
	}

}

func fixBson(s []byte) []byte {
	reg := regexp.MustCompile(`("_id")[ ]*[:][ ]*(ObjectId\(")[\w]{24}("\))[ ]*(,)[ ]*`)
	return reg.ReplaceAll(s, []byte(" "))
}

// Role ...
type Role []RoleElement

// RoleElement ...
type RoleElement struct {
	Name       string `json:"name"`
	Birthday   string `json:"birthday"`
	Age        string `json:"age"`
	Avatar     string `json:"avatar"`
	Height     string `json:"height"`
	Cup        Cup    `json:"cup"`
	Chest      string `json:"chest"`
	Waist      string `json:"waist"`
	Hipline    string `json:"hipline"`
	BirthPlace string `json:"birthPlace"`
	Hobby      string `json:"hobby"`
	Uncensored bool   `json:"uncensored"`
}

// Cup ...
type Cup string

const (
	// A ...
	A Cup = "A"
	// B ...
	B Cup = "B"
	// C ...
	C Cup = "C"
	// D ...
	D Cup = "D"
	// E ...
	E Cup = "E"
	// Empty ...
	Empty Cup = ""
	// F ...
	F Cup = "F"
	// G ...
	G Cup = "G"
	// H ...
	H Cup = "H"
	// I ...
	I Cup = "I"
	// J ...
	J Cup = "J"
	// K ...
	K Cup = "K"
	// M ...
	M Cup = "M"
	// O ...
	O Cup = "O"
)
