package main

import (
	"runtime"
	"syscall"
	"os/signal"
	"path/filepath"
	"os"
    "climbwall/utils"
    "climbwall/core"
)


func run_server (config * core.AppConfig) {
	core.Run_server_routine(config)
}


func wait_s(){
	var system_signal = make(chan os.Signal, 2)
	signal.Notify(system_signal, syscall.SIGINT, syscall.SIGHUP)
	for sig := range system_signal {
		if sig == syscall.SIGHUP || sig == syscall.SIGINT{
			utils.Logger.Printf("caught signal %v, exit", sig)
			os.Exit(0)
			
		} else {

			utils.Logger.Printf("XXX caught signal %v, exit", sig)
			os.Exit(0)
		}
	}
}


func main(){

	cpu := runtime.NumCPU()
	runtime.GOMAXPROCS(cpu)
	utils.Logger.Print("server log...")
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		utils.Logger.Print("can not combin config file path!")
		os.Exit(1)
	}
   
	config, err := core.Parse(dir + "/config.json")
	if err != nil {
		utils.Logger.Print("load config file failed!")
		os.Exit(1)
	}
	core.Dump_config(config)

	go run_server(config)

	wait_s()
}



