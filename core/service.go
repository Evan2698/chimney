package core

import (
	"net"
	"strconv"

	"github.com/Evan2698/chimney/utils"

	"github.com/Evan2698/chimney/config"
)

// Runclientsservice ...
func Runclientsservice(host string, app *config.AppConfig, p SocketService, f DataFlow, quit <-chan int) {
	all, err := net.Listen("tcp", host)
	if err != nil {
		utils.LOG.Print("local listen on   ip =", host, err)
		return
	}

	defer func() {
		utils.LOG.Println("listener will be close.^_^")
		all.Close()
		utils.LOG.Println("Runclientsservice is over!!!!^_^")
	}()

	for {
		someone, err := all.Accept()
		if err != nil {
			utils.LOG.Print("Accept failed: ", err)
			break
		}
		go handclientonesocket(someone, app, p, f)

		if quit != nil {
			select {
			case <-quit:
				utils.LOG.Println("will be exit!!")
				return
			default:
			}
		}
	}

	utils.LOG.Print("exit exit exit exit", err)
}

func handclientonesocket(o net.Conn, app *config.AppConfig, p SocketService, f DataFlow) {

	utils.SetSocketTimeout(o, app.Timeout)

	proxyhost := net.JoinHostPort(app.Server, strconv.Itoa(int(app.ServerPort)))
	con, err := createclientsocket(proxyhost, p)
	if err != nil {
		o.Close()
		utils.LOG.Print("create socket failed", err)
		return
	}
	utils.SetSocketTimeout(con, app.Timeout)

	ss := NewSocksSocket(con, app.Password, nil)
	proxysocket := NewSocketProxy(ss, app)
	h := NewSocksHandler(o, proxysocket, app)
	defer h.Close()

	err = h.Receive()
	if err != nil {
		utils.LOG.Print("client recv failed: ", err)
		return
	}
	h.Run(f)
}

// RunServerservice ..
func RunServerservice(host string, app *config.AppConfig, p SocketService, f DataFlow) {

	all, err := net.Listen("tcp", host)
	if err != nil {
		utils.LOG.Print("local listen on   ip =", host, err)
		return
	}
	defer all.Close()
	for {
		someone, err := all.Accept()
		if err != nil {
			utils.LOG.Print("remote socket failed to open", err)
			break
		}
		go handServeronesocket(someone, app, p, f)
	}
}

func handServeronesocket(o net.Conn, app *config.AppConfig, p SocketService, f DataFlow) {

	utils.SetSocketTimeout(o, app.Timeout)
	ss := NewSocksSocket(o, app.Password, nil)
	proxysocket := NewSocketProxy(ss, app)
	h := NewSocksHandler(nil, proxysocket, app)
	defer h.Close()

	err := h.Receive()
	if err != nil {
		utils.LOG.Print("client recv failed", err)
		return
	}
	h.Run(f)
}
