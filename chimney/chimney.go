package chimney

import (
	"github.com/Evan2698/chimney/core"
)

var gchimney *core.ListenHandle

// StartChimney ..
func StartChimney(s string, sport int, l string, lport int, pass string) bool {

	config := &core.AppConfig{
		ServerPort:   sport,
		LocalPort:    lport,
		LocalAddress: l,
		Server:       s,
		Password:     pass,
		Timeout:      600,
	}

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
