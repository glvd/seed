package seed

import (
	"github.com/ipfs/go-ipfs-api"
)

var sh = shell.NewShell("localhost:5001")

func AddDir(dir string) (s string, e error) {
	return sh.AddDir(dir)
}

func List(path string) (ls []*shell.LsLink, e error) {
	return sh.List(path)
}
