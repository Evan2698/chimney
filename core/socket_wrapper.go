//  +build !DRCLO

package core

import (
	"net"
	"strconv"
)

func Build_low_socket(ipString string, port int) (*CommonSocket, error) {

	host := net.JoinHostPort(ipString, strconv.Itoa(port))

	conn, err := net.Dial("tcp", host)

	return &CommonSocket{
		Remote: conn,
	}, err
}

func (con *CommonSocket) Close() {
	con.Remote.Close()
}
