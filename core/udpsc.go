package core

import (
	"net"
	"strconv"

	"chimney/security"

	"chimney/utils"

	"chimney/config"
)

//SclientRoutine for client
func SclientRoutine(app *config.AppConfig, p SocketService) (*net.UDPConn, error) {

	host := net.JoinHostPort(app.LocalAddress, strconv.Itoa(int(app.LocalPort)))

	udpaddr, err := net.ResolveUDPAddr("udp", host)
	if err != nil {
		utils.LOG.Println("can not resolve udp address", err)
		return nil, err
	}

	conn, err := net.ListenUDP("udp", udpaddr)
	if err != nil {
		utils.LOG.Println("can not resolve udp address", err)
		return nil, err
	}

	proxyhost := net.JoinHostPort(app.Server, strconv.Itoa(int(app.ServerPort)))

	go handles(conn, proxyhost, app.Password, p)

	return conn, nil

}

func handles(conn *net.UDPConn, proxy string, pw string, p SocketService) {
	defer conn.Close()

	for {
		buf := make([]byte, 4096)
		n, readdr, err := conn.ReadFromUDP(buf)
		if err != nil {
			utils.LOG.Println("udp read failed!!!!!")
			break
		}
		go handleoneudp(buf[:n], readdr, proxy, pw, conn, p)
	}

}

func handleoneudp(raw []byte, addr *net.UDPAddr, proxy string, pw string, root *net.UDPConn, p SocketService) {

	con, err := CreateCommonSocket(proxy, "udp", 60, p) //createclientsocket(proxy, p, "udp")
	if err != nil {
		utils.LOG.Print("can not connect udp server", proxy)
		return
	}
	defer con.Close()

	compressData, I, err := PackUDPData(pw, raw)
	if err != nil {
		utils.LOG.Print("package udp data error", err)
		return
	}

	_, err = con.Write(compressData)
	if err != nil {
		utils.LOG.Print("write udp server error!", err)
		return
	}
	var buf [4096]byte
	n, err := con.Read(buf[:])
	if err != nil {
		utils.LOG.Print("read from server error", err)
		return
	}

	out, err := I.Uncompress(buf[:n], security.MakeCompressKey(pw))
	if err != nil {
		utils.LOG.Print("uncompress failed", err)
		return
	}
	_, err = root.WriteToUDP(out, addr)

	utils.LOG.Print("write to client ", err)

}
