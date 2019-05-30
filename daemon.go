package seed

import (
	"gopkg.in/urfave/cli.v2"
	"io/ioutil"
	"path/filepath"
)

// DaemonCallback ...
type DaemonCallback func(path string)

// CmdDaemon ...
func CmdDaemon(flags []cli.Flag) *cli.Command {
	return &cli.Command{
		Name:    "daemon",
		Aliases: []string{"D"},
		Usage:   "daemon the filepath to process",
		Action: func(context *cli.Context) error {
			log.Info("daemon call")
			return nil
		},
		Subcommands: nil,
		Flags:       flags,
	}

	//for {
	//	if lists := procList(monitorPath, func(path string) {}); lists != nil {
	//		log.Info("processing", monitorPath)
	//
	//	}
	//	log.Info("waiting for new files")
	//	time.Sleep(30 * time.Second)
	//}
}

func procList(monitorPath string, dc DaemonCallback) (e error) {
	infos, e := ioutil.ReadDir(monitorPath)
	if e != nil {
		panic(e)
	}
	for _, info := range infos {
		name := info.Name()
		log.Infof("name:%s ext:%s", name, filepath.Ext(name))
		if !info.IsDir() && filepath.Ext(name) != "filepart" {
			fullpath := filepath.Join(monitorPath, info.Name())
			dc(fullpath)
		}
		if name == "success" {
			continue
		}
	}
	return nil
}
