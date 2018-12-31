package chimney

import (
	"github.com/Evan2698/chimney/core"
)

var gchimney *core.ListenHandle

//ISocket ...
type ISocket core.SocketService

// IDataFlow ..
type IDataFlow core.DataFlowService

//Register ..
func Register(v ISocket, k IDataFlow) {
	core.GFlow = k
	core.GSocketInterface = v
}

// StartChimney ..
func StartChimney(s string,
	sport int,
	l string,
	lport int,
	pass string,
	path string) bool {

	config := &core.AppConfig{
		ServerPort:   sport,
		LocalPort:    lport,
		LocalAddress: l,
		Server:       s,
		Password:     pass,
		Timeout:      1000,
	}

	core.GUNIXPATH = path

	gchimney = core.RunAndroidListenLoop(config)

	return (gchimney != nil)
}

// StopChimney ..
func StopChimney() bool {

	if gchimney != nil {
		core.StopAndroidWorld(gchimney)
	}

	return true

}
