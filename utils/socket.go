package utils

import (
	"net"
	"time"
)

//SetSocketTimeout ..
func SetSocketTimeout(con net.Conn, tm int) {
	if con != nil && tm != 0 {
		readTimeout := time.Duration(tm) * time.Second
		v := time.Now().Add(readTimeout)
		con.SetReadDeadline(v)
		con.SetWriteDeadline(v)
		con.SetDeadline(v)
	}
}
