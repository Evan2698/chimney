package chimney

import (
	"github.com/eycorsican/go-tun2socks/cmd/ago"
)

//StartNetstackService ..
func StartNetstackService(fd int, socks string, dns string) {

	ago.Tun2SocksMain(fd)
}

//StopNetStackService ..
func StopNetStackService() {

	ago.StopTun2SocksMain()
}
