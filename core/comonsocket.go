package core

import (
	"net"
)

type CommonSocket struct {
	Remote net.Conn
	Fd     uintptr
}
