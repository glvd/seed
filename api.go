package seed

import (
	"github.com/godcong/go-ipfs-restapi"
	"github.com/sirupsen/logrus"
)

var rest *shell.Shell

// InitShell ...
func InitShell(s string) {
	logrus.Info("ipfs shell:", s)
	rest = shell.NewShell(s)
}
