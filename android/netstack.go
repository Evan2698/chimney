package chimney

import "github.com/Evan2698/android-netstack/cmd/astack"

//StartNetstackService ..
func StartNetstackService(fd int, socks string, dns string) {
	astack.StartNetstackService(fd, socks, dns)
	//mobile.StartService(fd, socks, dns)

}

//StopNetStackService ..
func StopNetStackService() {
	astack.StopNetStackService()
	//mobile.StopService()
}
