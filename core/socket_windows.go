package core

import (
	"net"
	"strconv"

	"github.com/Evan2698/chimney/utils"

	"github.com/Evan2698/chimney/config"
)

func createclientsocket(p SocketService, network string, app *config.AppConfig) (net.Conn, error) {
	host := net.JoinHostPort(app.Server, strconv.Itoa(int(app.ServerPort)))
	outcon, err := CreateCommonSocket(host, network, app.Timeout, p)
	utils.LOG.Print("create as common socket!!")
	return outcon, err
}

func CreateCommonSocket(host string, network string, timeout uint32, p SocketService) (net.Conn, error) {
	outcon, err := net.Dial(network, host)
	if err == nil {
		utils.SetSocketTimeout(outcon, timeout)
	}

	return outcon, err
}
