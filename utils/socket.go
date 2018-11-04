package utils

import (
	"net"
	"time"
)

// SetReadTimeOut function
//
func SetReadTimeOut(con net.Conn, timeout int) {
	if con != nil && timeout != 0 {
		readTimeout := time.Duration(timeout) * time.Second
		con.SetReadDeadline(time.Now().Add(readTimeout))
	}

}
