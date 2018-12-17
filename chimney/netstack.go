package chimney

import (
	"io"
	"os"

	netstack "github.com/Evan2698/android-netstack"
)

var netdevice *netstack.NetStack2Socks

//StartNetstackService ..
func StartNetstackService(fd int, socks string, dns string) {

	f := os.NewFile((uintptr)(fd), "")
	var rwc io.ReadWriteCloser
	rwc = f
	dnsArray := make([]string, 1)
	dnsArray[0] = dns
	netdevice = netstack.New(rwc, socks, dnsArray, true, true)

	go func(device *netstack.NetStack2Socks) {
		netdevice.Run()

	}(netdevice)
}

//StopNetStackService ..
func StopNetStackService() {

	if netdevice != nil {
		netdevice.Stop()
	}
	netdevice = nil
}
