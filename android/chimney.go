package chimney

import (
	"socks5/config"
	"socks5/core"
)

//ISocket ...
type ISocket core.SocketService

// IDataFlow ..
type IDataFlow core.DataFlow

//Register ..
func Register(v ISocket, k IDataFlow) {
	//core.GFlow = k
	//core.GSocketInterface = v
}

// StartChimney ..
func StartChimney(s string,
	sport int,
	l string,
	lport int,
	pass string,
	path string) bool {

	config := &config.AppConfig{
		ServerPort:   sport,
		LocalPort:    lport,
		LocalAddress: l,
		Server:       s,
		Password:     pass,
		Timeout:      1000,
	}

	return config == nil
}

// StopChimney ..
func StopChimney() bool {

	//if gchimney != nil {
	//	core.StopAndroidWorld(gchimney)
	//}

	return true

}
