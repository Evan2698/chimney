package mobile

import (
	"github.com/Evan2698/chimney/lwip2socks/common/dns/cache"
	"github.com/Evan2698/chimney/lwip2socks/core"
	"github.com/Evan2698/chimney/lwip2socks/proxy/socks"
	"io"
	"log"
	"os"
	"time"
)

var dnsCache = cache.NewDNSCache()

var lwipWriter = core.NewLWIPStack()

var tun *os.File

// StartService ...
func StartService(fd int, proxy string, dns string) bool {
	tun = os.NewFile(uintptr(fd), "")

	core.RegisterTCPConnHandler(socks.NewTCPHandler("127.0.0.1", 1080))
	core.RegisterUDPConnHandler(socks.NewUDPHandler("127.0.0.1", 1080, 180*time.Second, dnsCache))

	core.RegisterOutputFn(func(data []byte) (int, error) {
		return tun.Write(data)
	})

	go func() {
		n, err := io.Copy(lwipWriter, tun)
		if err != nil {
			log.Println("tun will exit!!!", err)
		}
		log.Println("log failed.", n)
	}()

	return true

}

// StopService ...
func StopService() {
	if tun != nil {
		tun.Close()
		tun = nil
	}
	time.Sleep(4 * time.Second)
}
