package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strconv"
	"syscall"

	"github.com/Evan2698/chimney/utils"

	"github.com/Evan2698/chimney/core"

	"github.com/Evan2698/chimney/config"
)

func main() {
	var configpath string
	cpu := runtime.NumCPU()
	runtime.GOMAXPROCS(cpu * 4)
	utils.LOG.Print("I AM SERVER!!!!!!")

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		utils.LOG.Print("can not combin config file path!")
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

	config, err := config.Parse(configpath)
	if err != nil {
		utils.LOG.Print("load config file failed!", err)
		fmt.Println("load config file failed!", err)
		os.Exit(1)
	}

	host := net.JoinHostPort(config.Server, strconv.Itoa(int(config.ServerPort)))

	go core.RunServerservice(host, config, nil, nil)

	if !config.UseQuic {

		con, err := core.SServerRoutine(config)
		if err != nil {
			utils.LOG.Println("udp server failed", err)
			os.Exit(1)
		}
		defer con.Close()
	}

	waits()
}

func waits() {
	var systemsignal = make(chan os.Signal, 2)
	signal.Notify(systemsignal, syscall.SIGINT, syscall.SIGHUP)
	for sig := range systemsignal {
		if sig == syscall.SIGHUP || sig == syscall.SIGINT {
			utils.LOG.Printf("caught signal %v, exit", sig)
			os.Exit(0)

		} else {

			utils.LOG.Printf("XXX caught signal %v, exit", sig)
			os.Exit(0)
		}
	}
}
