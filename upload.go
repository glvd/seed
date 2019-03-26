package seed

import (
	log "github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
)

func prefix(s string) (ret string) {
	ret = "/ipfs/" + s
	return
}

func Upload(source *VideoSource) (e error) {
	if source == nil {
		return xerrors.New("nil source")
	}

	s, e := AddDir(source.FilePath)
	if e != nil {
		return e
	}
	ls, e := List(prefix(s))
	for _, v := range ls {
		log.Info(v)
	}
	return nil
}

func Load(path string) []*VideoSource {
	var vs []*VideoSource

	e := ReadJSON(path, &vs)
	if e != nil {
		return nil
	}
	return vs
}
