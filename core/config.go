package core

import (
	"time"
	"encoding/json"
	"io/ioutil"
	"os"

)


type AppConfig struct{
	ServerPort int         `json:"server_port"`
	LocalPort  int         `json:"local_port"`
	Password   string      `json:"password"`
	Timeout    int         `json:"timeout"`
	Server     string      `json:"server"`
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
	
	config.Timeout = int (time.Duration(config.Timeout) * time.Second)
	
	return config, nil
}