package main

import (
	"climbwall/core"
	"climbwall/utils"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
)

func runlocal(config *core.AppConfig) {
	core.Run_Local_routine(config)
}

func waits() {
	var systemsignal = make(chan os.Signal, 2)
	signal.Notify(systemsignal, syscall.SIGINT, syscall.SIGHUP)
	for sig := range systemsignal {
		if sig == syscall.SIGHUP || sig == syscall.SIGINT {
			utils.Logger.Printf("caught signal %v, exit", sig)
			os.Exit(0)

		} else {

			utils.Logger.Printf("XXX caught signal %v, exit", sig)
			os.Exit(0)
		}
	}
}

func main() {

	var configpath string
	cpu := runtime.NumCPU()
	runtime.GOMAXPROCS(cpu)
	utils.Logger.Print("local log.....")
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		utils.Logger.Print("can not combin config file path!")
		os.Exit(1)
	}

	if len(os.Args) > 1 {
		co := flag.String("c", "", "please input config file")
		flag.Parse()
		configpath = *co
	} else {
		configpath = dir + "/config.json"
	}

	if (len(configpath)) == 0 {
		fmt.Println("config file path is incorrect!!", configpath)
		os.Exit(1)
	}

	config, err := core.Parse(configpath)
	if err != nil {
		utils.Logger.Print("load config file failed!", err)
		fmt.Println("load config file failed!", err)
		os.Exit(1)
	}

	core.Dump_config(config)

	go runlocal(config)

	waits()
}
