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
		Name:        "shell",
		Aliases:     []string{"s"},
		Value:       "localhost:5001",
		Usage:       "set the ipfs api address:port",
		DefaultText: "localhost:5001",
	}

	database := &cli.StringFlag{
		Name:    "database",
		Aliases: []string{"db"},
		Usage:   "set the database path",
	}

	config := &cli.StringFlag{
		Name:    "config",
		Aliases: []string{"c"},
		Usage:   "sometime need set config",
	}

	userpass := &cli.StringFlag{
		Name:    "userpass",
		Aliases: []string{"u"},
		Usage:   "set the database user:password",
	}

	json := &cli.StringFlag{
		Name:    "json",
		Aliases: []string{"j"},
		Usage:   "set the json file path",
	}

	boot := &cli.StringFlag{
		Name:    "bootstrap",
		Aliases: []string{"b"},
		Usage:   "set the ipfs bootstrap swarm address to quick connect",
	}

	pin := &cli.BoolFlag{
		Name:    "pin",
		Aliases: []string{"p"},
		Usage:   "set to pin on ipfs",
	}

	sync := &cli.BoolFlag{
		Name:    "sync",
		Aliases: []string{"n"},
		Usage:   "check if the video is synced",
	}

	return []cli.Flag{
		shell, database, userpass, json, boot, pin, sync, config,
	}

}

func runApp() error {
	flags := globalFlags()
	app := &cli.App{
		Version: "v0.0.1",
		Name:    "seed",
		Usage:   "seed is a video manage tool use ipfs,eth,sqlite3 and so on.",
		Action: func(c *cli.Context) error {
			log.Info(c.String("s"))
			return nil
		},
		Flags: flags,
	}

	app.Commands = []*cli.Command{
		seed.CmdDaemon(app),
		seed.CmdProcess(app),
		seed.CmdContract(app),
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
