package main

import (
	"fmt"
	"github.com/godcong/go-trait"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
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

func main() {

	app := &cli.App{
		Version: "v0.0.1",
		Name:    "seed",
		Usage:   "seed is a video manage tool use ipfs,eth,sqlite3,and so on.",
		Action: func(c *cli.Context) error {
			log.Info(c.String("s"))
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:    "complete",
				Aliases: []string{"c"},
				Usage:   "complete a task on the list",
				Action: func(c *cli.Context) error {
					return nil
				},
			},
			{
				Name:  "process",
				Usage: "add a task to the list",
				Action: func(c *cli.Context) error {
					return nil
				},
			},
		},
	}

	//cli.HelpFlag = &cli.BoolFlag{
	//	Name:    "help",
	//	Aliases: []string{"h"},
	//	Usage:   "show help",
	//	Value:   true,
	//}
	shell := &cli.StringFlag{
		Name:    "shell",
		Aliases: []string{"s"},
		Value:   "localhost:5001",
		Usage:   "set the ipfs api address:port",
		//DefaultText: "localhost:5001",
	}

	quick := &cli.BoolFlag{
		Name:    "quick",
		Aliases: []string{"q"},
		Usage:   "set the ipfs api address:port",
		//DefaultText: "localhost:5001",
	}

	app.Flags = []cli.Flag{
		shell, quick,
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
	//defer func() {
	//	if e := recover(); e != nil {
	//		log.Panic(e)
	//	}
	//}()
	//
	////shell := rootCmd.PersistentFlags().StringP("shell", "s", "localhost:5001", "set the ipfs api address:port")
	//path := rootCmd.PersistentFlags().StringP("path", "p", "seed.json", "set the path to load video source or source path")
	//config := rootCmd.PersistentFlags().StringP("config", "c", "config.toml", "file config")
	//tipe := rootCmd.PersistentFlags().StringP("type", "t", "json", "transfer to types")
	//
	//quick := rootCmd.PersistentFlags().BoolP("quick", "q", false, "process with only filepath,no detail")
	//poster := rootCmd.PersistentFlags().BoolP("poster", "o", false, "only pin poster")
	//check := rootCmd.PersistentFlags().BoolP("check", "k", true, "check if the video is synced")
	//swarm := rootCmd.PersistentFlags().StringP("swarm", "w", bootIPFS, "quick connect to ipfs")
	//pin := rootCmd.PersistentFlags().BoolP("pin", "i", false, "check if need pin")
	//
	//var cmdDaemon = &cobra.Command{
	//	Use:   "daemon",
	//	Short: "daemon the path",
	//	Long:  "daemon the filepath to process",
	//	Run: func(cmd *cobra.Command, args []string) {
	//		seed.DaemonStart(*path)
	//	},
	//}
	//var cmdContract = &cobra.Command{
	//	Use:   "contract",
	//	Short: "contract the data to eth.",
	//	Long:  `contract the data information from database to eth contract`,
	//	Run: func(cmd *cobra.Command, args []string) {
	//		key := ""
	//		if len(args) > 0 {
	//			key = args[0]
	//		}
	//		if err := seed.Contract(key); err != nil {
	//			panic(err)
	//		}
	//	},
	//}
	//
	//var cmdVerify = &cobra.Command{
	//	Use:   "verify",
	//	Short: "verify the json file",
	//	Long:  `verify the json file is correct before transfer`,
	//	Run: func(cmd *cobra.Command, args []string) {
	//		vs := seed.Load(*path)
	//		for _, v := range vs {
	//			e := seed.Verify(v)
	//			if e != nil {
	//				panic(e)
	//			}
	//			log.Infof("%+v", v)
	//		}
	//		log.Info("success")
	//	},
	//}
	//
	//var cmdProcess = &cobra.Command{
	//	Use:   "process",
	//	Short: "process split and upload to ipfs",
	//	Long: `process split a video file to hls ,then upload to ipfs.
	//	after that return a ipfs hash info json.`,
	//	//Args: cobra.MinimumNArgs(1),
	//	Run: func(cmd *cobra.Command, args []string) {
	//		seed.InitShell(*shell)
	//		e := model.InitDB()
	//		if e != nil {
	//			panic(e)
	//		}
	//
	//		if !*quick {
	//			vs := seed.Load(*path)
	//			for _, v := range vs {
	//				e := seed.Process(v)
	//				if e != nil {
	//					panic(e)
	//				}
	//				log.Infof("%+v", v)
	//			}
	//			return
	//		}
	//		if err := seed.QuickProcess(*path, *pin); err != nil {
	//			log.Panic(err)
	//			return
	//		}
	//
	//	},
	//}
	//
	//var cmdUpdate = &cobra.Command{
	//	Use:   "update",
	//	Short: "update the information",
	//	Long:  `update only update the video information to new information`,
	//	//Args: cobra.MinimumNArgs(1),
	//	Run: func(cmd *cobra.Command, args []string) {
	//		fmt.Println(args)
	//		if len(args) < 1 {
	//			fmt.Println("nothing process")
	//			return
	//		}
	//		seed.InitShell(*shell)
	//		e := model.InitDB()
	//		if e != nil {
	//			panic(e)
	//		}
	//		vs := seed.Load(*path)
	//		for _, v := range vs {
	//			e := seed.Update(args[0], v)
	//			if e != nil {
	//				log.Error(e)
	//				continue
	//			}
	//		}
	//
	//		log.Infof("%+v", vs[0])
	//	},
	//}
	//
	//var cmdPin = &cobra.Command{
	//	Use:   "pin",
	//	Short: "pin the video to ipfs",
	//	Long:  `pin the video to ipfs, then user can get it more quickly`,
	//	//Args: cobra.MinimumNArgs(1),
	//	Run: func(cmd *cobra.Command, args []string) {
	//		pin := ""
	//		if len(args) >= 1 {
	//			pin = args[0]
	//		}
	//		seed.InitShell(*shell)
	//		e := model.InitDB()
	//		if e != nil {
	//			panic(e)
	//		}
	//		go seed.PoolSwarmConnect()
	//		seed.SwarmAddAddress(*swarm)
	//		if !*quick {
	//			e = seed.Pin(pin, *poster, *check)
	//			if e != nil {
	//				panic(e)
	//			}
	//			return
	//		}
	//
	//		if err := seed.QuickPin(pin, *check); err != nil {
	//			return
	//		}
	//
	//	},
	//}
	//
	//var cmdTransfer = &cobra.Command{
	//	Use:   "transfer",
	//	Short: "transfer video data info to json",
	//	Long:  `transfer will output a json file from video info db.`,
	//	//Args:  cobra.MinimumNArgs(1),
	//	Run: func(cmd *cobra.Command, args []string) {
	//		seed.InitShell(*shell)
	//		e := model.InitDB()
	//		if e != nil {
	//			log.Panic(e)
	//		}
	//
	//		switch *tipe {
	//		case "mysql":
	//			eng, e := model.InitSync(*tipe, *config)
	//			if e != nil {
	//				panic(e)
	//			}
	//			e = seed.TransferTo(eng, 1)
	//			if e != nil {
	//				panic(e)
	//			}
	//		case "sqlite3":
	//			eng, e := model.InitSync(*tipe, *path)
	//			if e != nil {
	//				panic(e)
	//			}
	//			e = seed.TransferTo(eng, 1)
	//			if e != nil {
	//				panic(e)
	//			}
	//		default:
	//			if err := seed.Transfer(); err != nil {
	//				panic(e)
	//			}
	//		}
	//	},
	//}
	//rootCmd.AddCommand(cmdProcess, cmdTransfer, cmdUpdate, cmdPin, cmdVerify, cmdContract, cmdDaemon)
	//rootCmd.SuggestionsMinimumDistance = 1
	//Execute()
}

// Execute ...
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
