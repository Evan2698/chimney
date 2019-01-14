package core

import (
	"github.com/Evan2698/chimney/security"
)

// SocketService ...
type SocketService interface {
	Protect(fd int) bool
}

// DataFlow ...
type DataFlow interface {
	Update(up, down int64)
}

// SSocket ..
type SSocket interface {
	Read() (buf []byte, err error)
	Write(p []byte) error
	Close() error
	SetI(i security.EncryptThings)
}

// SocksProxy ...
type SocksProxy interface {
	Read() (buf []byte, err error)
	Write(p []byte) error
	SetEncrypt(I security.EncryptThings)
	Close() error
	Connect(remoteaddr []byte) error
}

// SocksHandler ...
type SocksHandler interface {
	Close() error
	Receive(p SocketService) error
	Run(f DataFlow)
}
