package core

import (
	"net"
)

func createclientsocket(host string, p SocketService) (net.Conn, error) {
	con, err := net.Dial("tcp", host)
	return con, err
}
