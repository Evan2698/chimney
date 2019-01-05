package core

import (
	"net"
	"socks5/utils"
)

func createclientsocket(host string, p SocketService) (net.Conn, error) {

	con, err := net.Dial("tcp", host)
	if err != nil {
		return nil, err
	}

	if p != nil {
		tcp, ok := con.(*net.TCPConn)
		if ok {
			f, err := tcp.File()
			if err == nil {
				fd := f.Fd()
				p.Protect(int(fd))
				f.Close()
			} else {
				utils.LOG.Print("can not get file descriptor,", err)
			}
		}
	}

	return con, nil

}
