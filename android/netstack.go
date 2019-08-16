package chimney

import "lwip2socks/mobile"

//StartNetstackService ..
func StartNetstackService(fd int, socks string, dns string) {
	//astack.StartNetstackService(fd, socks, dns)
	//mobile.StartService(fd, socks, dns)
	mobile.StartService(fd, socks, dns)

}

//StopNetStackService ..
func StopNetStackService() {
	//astack.StopNetStackService()
	mobile.StopService()
}
