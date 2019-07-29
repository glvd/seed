package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

var log = trait.NewZapSugar()

func main() {
	args := os.Args
	dir, err := os.Getwd()
	if len(args) > 1 {
		err = nil
		dir = args[1]
	}
	if err != nil {
		log.Info("wd:", err)
		return
	}

	var role Role
	err = json.Unmarshal(&role)
	if err != nil {
		panic(err)
	}

	for _, ele := range ([]RoleElement)(role) {
		ext := filepath.Ext(ele.Avatar)
		log.With("from", ele.Avatar, "to", ele.Name+ext).Info("rename")
		_ = os.Rename(ele.Avatar, ele.Name+ext)
	}

}

type Role []RoleElement

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

type Cup string

const (
	A     Cup = "A"
	B     Cup = "B"
	C     Cup = "C"
	D     Cup = "D"
	E     Cup = "E"
	Empty Cup = ""
	F     Cup = "F"
	G     Cup = "G"
	H     Cup = "H"
	I     Cup = "I"
	J     Cup = "J"
	K     Cup = "K"
	M     Cup = "M"
	O     Cup = "O"
)
