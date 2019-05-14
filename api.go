package seed

import (
	"github.com/godcong/go-ipfs-restapi"
	"github.com/sirupsen/logrus"
	"time"
)

var rest *shell.Shell

// InitShell ...
func InitShell(s string) {
	logrus.Info("ipfs shell:", s)
	rest = shell.NewShell(s)
}

// QuickConnect ...
func QuickConnect(addr string) {
	var e error
	go func() {
		for {
			e = SwarmConnect(addr)
			if e != nil {
				return
			}
			time.Sleep(30 * time.Second)
		}
	}()
}
