package config

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/Evan2698/chimney/utils"
)

// AppConfig ..
type AppConfig struct {
	ServerPort   uint16 `json:"server_port"`
	LocalPort    uint16 `json:"local_port"`
	LocalAddress string `json:"local_address"`
	Password     string `json:"password"`
	Timeout      uint32 `json:"timeout"`
	Server       string `json:"server"`
	SSLRaw       bool   `json:"sslraw"`
	UseQuic      bool   `json:"quic"`
	QuicPort     uint16 `json:"quic_port"`
}

// Parse ..
func Parse(path string) (config *AppConfig, err error) {
	file, err := os.Open(path) // For read access.
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	config = &AppConfig{}
	if err = json.Unmarshal(data, config); err != nil {
		return nil, err
	}

	return config, nil
}

// DumpConfig ..
func DumpConfig(config *AppConfig) {
	utils.LOG.Print("server :", config.Server)
	utils.LOG.Print("server_port :", config.ServerPort)
	utils.LOG.Print("local_port :", config.LocalPort)
	utils.LOG.Print("password :", config.Password)
	utils.LOG.Print("timeout :", config.Timeout)
	utils.LOG.Print("sslraw :", config.SSLRaw)
	utils.LOG.Print("UseQuic :", config.UseQuic)
	utils.LOG.Print("QuicPort :", config.QuicPort)
}
