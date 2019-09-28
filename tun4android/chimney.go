package tun4android

import (
	"net"
	"time"

	"github.com/Evan2698/chimney/utils"

	"github.com/Evan2698/chimney/core"

	"github.com/Evan2698/chimney/config"
)

//ISocket ...
type ISocket core.SocketService

// IDataFlow ..
type IDataFlow core.DataFlow

var flow IDataFlow
var sockets ISocket

var quit chan int32

var udpconn *net.UDPConn

//Register ..
func Register(v ISocket, k IDataFlow) {
	flow = k
	sockets = v
}

// StartChimney ..
func StartChimney(s string,
	sport uint16,
	l string,
	lport uint16,
	pass string,
	path string) bool {

	config := &config.AppConfig{
		ServerPort:   sport,
		LocalPort:    lport,
		LocalAddress: l,
		Server:       s,
		Password:     pass,
		Timeout:      30,
		SSLRaw:       true,
		QuicPort:     443,
		UseQuic:      false,
	}
	quit = make(chan int32, 1)
	go core.Runclientsservice("127.0.0.1:1080", config, sockets, flow, quit)

	var err error
	go func() {
		udpconn, err = core.SclientRoutine(config, sockets)
		if err != nil {
			utils.LOG.Println("ERROR! ERROR!")
		}

	}()

	return true
}

// StopChimney ..
func StopChimney() bool {

	if udpconn != nil {
		udpconn.Close()
		udpconn = nil
	}

	if quit != nil {
		utils.LOG.Println("stop stop stop stop stop")
		quit <- 1
		con, _ := net.Dial("tcp", "127.0.0.1:1080")
		con.Close()
		utils.LOG.Println("end end end end stop")
		time.Sleep(time.Second * 2)
		close(quit)
		quit = nil
	}

	return true
}
