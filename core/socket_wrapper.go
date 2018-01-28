//  +build !DRCLO

package core

import (
	"net"
	"strconv"
)

func Build_low_socket(ipString string, port int) (net.Conn, int, error) {

	host := net.JoinHostPort(ipString, strconv.Itoa(port))

	conn, err := net.Dial("tcp", host)

	return conn, -1, err
}
