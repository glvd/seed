package main

import (
	"fmt"
	"github.com/godcong/go-trait"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/yinhevr/seed"
	"github.com/yinhevr/seed/model"
	"os"
)

var rootCmd = &cobra.Command{
	Use:        "seed",
	Aliases:    nil,
	SuggestFor: nil,
	Short:      "Seed is a ipfs video split upload tool",
	Long: `Seed make you split and upload video file with hls in one step
				Complete documentation is available at https://github.com/yinhevr/seed`,
}

func main() {
	defer func() {
		if e := recover(); e != nil {
			log.Panic(e)
		}
	}()
	path := rootCmd.PersistentFlags().StringP("path", "p", "seed.json", "set the path to load video source")
	shell := rootCmd.PersistentFlags().StringP("shell", "s", "localhost:5001", "set the ipfs api port")
	//action := rootCmd.PersistentFlags().StringP("action", "a", "cmdProcess", "set action to do something")
	quick := rootCmd.PersistentFlags().BoolP("quick", "q", false, "process with only filepath,no detail")
	config := rootCmd.PersistentFlags().StringP("config", "c", "config.toml", "file config")
	tipe := rootCmd.PersistentFlags().StringP("type", "t", "json", "transfer to types")

	check := rootCmd.PersistentFlags().BoolP("check", "k", false, "check if the video is synced")
	trait.InitRotateLog("logs/seed.log", trait.RotateLogLevel(trait.RotateLogDebug))

	var cmdContract = &cobra.Command{
		Use:   "contract",
		Short: "contract the data to eth.",
		Long:  `contract the data information from database to eth contract`,
		Run: func(cmd *cobra.Command, args []string) {
			key := ""
			if len(args) > 0 {
				key = args[0]
			}
			if err := seed.Contract(key); err != nil {
				panic(err)
			}
		},
	}

	var cmdVerify = &cobra.Command{
		Use:   "verify",
		Short: "verify the json file",
		Long:  `verify the json file is correct before transfer`,
		Run: func(cmd *cobra.Command, args []string) {
			vs := seed.Load(*path)
			for _, v := range vs {
				e := seed.Verify(v)
				if e != nil {
					log.Panic(e)
				}
				log.Infof("%+v", v)
			}
			log.Info("success")
		},
	}

	var cmdProcess = &cobra.Command{
		Use:   "process",
		Short: "process split and upload to ipfs",
		Long: `process split a video file to hls ,then upload to ipfs. 
		after that return a ipfs hash info json.`,
		//Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			seed.InitShell(*shell)
			e := model.InitDB()
			if e != nil {
				log.Panic(e)
			}
			if !*quick {
				vs := seed.Load(*path)
				for _, v := range vs {
					e := seed.Process(v)
					if e != nil {
						log.Panic(e)
					}
					log.Infof("%+v", v)
				}
				return
			}
			if err := seed.QuickProcess(*path); err != nil {
				log.Panic(err)
				return
			}

		},
	}

	var cmdUpdate = &cobra.Command{
		Use:   "update",
		Short: "update the information",
		Long:  `update only update the video information to new information`,
		//Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(args)
			if len(args) < 1 {
				fmt.Println("nothing process")
				return
			}
			seed.InitShell(*shell)
			e := model.InitDB()
			if e != nil {
				log.Panic(e)
			}
			vs := seed.Load(*path)
			for _, v := range vs {
				e := seed.Update(args[0], v)
				if e != nil {
					log.Error(e)
					continue
				}
			}

			log.Infof("%+v", vs[0])
		},
	}

	var cmdPin = &cobra.Command{
		Use:   "pin",
		Short: "pin the video to ipfs",
		Long:  `pin the video to ipfs, then user can get it more quickly`,
		//Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			pin := ""
			if len(args) >= 1 {
				pin = args[0]
			}
			seed.InitShell(*shell)
			e := model.InitDB()
			if e != nil {
				log.Panic(e)
			}
			if !*quick {
				e = seed.Pin(pin, *check)
				if e != nil {
					log.Panic(e)
				}
				return
			}

			if err := seed.QuickPin(pin, *check); err != nil {
				return
			}

		},
	}

	var cmdTransfer = &cobra.Command{
		Use:   "transfer",
		Short: "transfer video data info to json",
		Long:  `transfer will output a json file from video info db.`,
		//Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			seed.InitShell(*shell)
			e := model.InitDB()
			if e != nil {
				log.Panic(e)
			}
			if *tipe == "mysql" {
				eng, e := model.InitSync(*config)
				if e != nil {
					panic(e)
				}
				e = seed.TransferMysql(eng, 1)
				if e != nil {
					panic(e)
				}
				return
			}
			if err := seed.Transfer(); err != nil {
				log.Panic(e)
			}
		},
	}
	rootCmd.AddCommand(cmdProcess, cmdTransfer, cmdUpdate, cmdPin, cmdVerify, cmdContract)
	rootCmd.SuggestionsMinimumDistance = 1
	Execute()
}

// Execute ...
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
