package core

import (
	"github.com/Evan2698/climbwall/utils"
	"encoding/json"
	"io/ioutil"
	"os"
)

type AppConfig struct {
	ServerPort int    `json:"server_port"`
	LocalPort  int    `json:"local_port"`
	Password   string `json:"password"`
	Timeout    int    `json:"timeout"`
	Server     string `json:"server"`
}

func Parse(path string) (config *AppConfig, err error) {
	file, err := os.Open(path) // For read access.
	if err != nil {
		return
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return
	}

	config = &AppConfig{}
	if err = json.Unmarshal(data, config); err != nil {
		return nil, err
	}

	return config, nil
}

func Dump_config(config *AppConfig) {
	utils.Logger.Print("server :", config.Server)
	utils.Logger.Print("server_port :", config.ServerPort)
	utils.Logger.Print("local_port :", config.LocalPort)
	utils.Logger.Print("password :", config.Password)
	utils.Logger.Print("timeout :", config.Timeout)
}
