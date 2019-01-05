package chimney

import (
	"net"

	"github.com/Evan2698/chimney/core"

	"github.com/Evan2698/chimney/config"
)

//ISocket ...
type ISocket core.SocketService

// IDataFlow ..
type IDataFlow core.DataFlow

var flow IDataFlow
var sockets ISocket
var handle net.Listener

//Register ..
func Register(v ISocket, k IDataFlow) {
	flow = k
	sockets = v
}

// StartChimney ..
func StartChimney(s string,
	sport int,
	l string,
	lport int,
	pass string,
	path string) bool {

	ch := make(chan net.Listener, 1)

	config := &config.AppConfig{
		ServerPort:   sport,
		LocalPort:    lport,
		LocalAddress: l,
		Server:       s,
		Password:     pass,
		Timeout:      600,
	}

	go core.Runclientsservice("127.0.0.1:1080", config, sockets, flow, ch)
	handle, _ := <-ch

	close(ch)
	return handle != nil
}

// StopChimney ..
func StopChimney() bool {

	if handle != nil {
		handle.Close()
		handle = nil
	}
	return true

}
