package main

import (
	"fmt"
	"github.com/godcong/go-trait"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
	"github.com/yinhevr/seed"
	"gopkg.in/urfave/cli.v2"
	"os"
	"sort"
)

const bootIPFS = "/ip4/47.101.169.94/tcp/4001/ipfs/QmeF1HVnBYTzFFLGm4VmAsHM4M7zZS3WUYx62PiKC2sqRq"

var log = trait.NewZapSugar()

var rootCmd = &cobra.Command{
	Use:        "seed",
	Aliases:    nil,
	SuggestFor: nil,
	Short:      "Seed is a ipfs video split upload tool",
	Long: `Seed make you split and upload video file with hls in one step
				Complete documentation is available at https://github.com/yinhevr/seed`,
}

func globalFlags() []cli.Flag {
	shell := &cli.StringFlag{
		Name:    "shell",
		Aliases: []string{"s"},
		//Value:       "localhost:5001",
		Usage:       "set the ipfs api address:port",
		DefaultText: "localhost:5001",
	}

	quick := &cli.BoolFlag{
		Name:    "quick",
		Aliases: []string{"q"},
		Usage:   "set the ipfs api address:port",
	}
	return []cli.Flag{
		shell, quick,
	}

}

func runApp() error {
	flags := globalFlags()
	app := &cli.App{
		Version: "v0.0.1",
		Name:    "seed",
		Usage:   "seed is a video manage tool use ipfs,eth,sqlite3,and so on.",
		Action: func(c *cli.Context) error {
			log.Info(c.String("s"))
			return nil
		},
		Commands: []*cli.Command{
			seed.CmdDaemon(flags),
			{
				Name:  "process",
				Usage: "add a task to the list",
				Action: func(c *cli.Context) error {
					return nil
				},
			},
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))
	err := app.Run(os.Args)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	if err := runApp(); err != nil {
		panic(err)
	}
	os.Exit(0)
}

// Execute ...
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
