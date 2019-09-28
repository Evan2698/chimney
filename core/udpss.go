package core

import (
	"net"
	"strconv"

	"github.com/Evan2698/chimney/security"

	"github.com/Evan2698/chimney/utils"

	"github.com/Evan2698/chimney/config"
)

func SServerRoutine(app *config.AppConfig) (*net.UDPConn, error) {

	host := net.JoinHostPort(app.Server, strconv.Itoa(int(app.ServerPort)))

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

	go handlesServer(conn, app.Password)

	return conn, nil

}

func handlesServer(conn *net.UDPConn, pw string) {
	defer conn.Close()

	for {
		buf := make([]byte, 4096)
		n, readdr, err := conn.ReadFromUDP(buf)
		if err != nil {
			utils.LOG.Println("udp read failed!!!!!")
			break
		}
		go handleoneudps(buf[:n], readdr, pw, conn)
	}

}

func handleoneudps(raw []byte, addr *net.UDPAddr, pw string, root *net.UDPConn) {

	o, dest, I, err := TryparseUDPProtocol(raw, pw)
	if err != nil {
		utils.LOG.Print("can not parse udp package", err)
		return
	}

	con, err := CreateCommonSocket(dest, "udp", 60, nil) //createclientsocket(dest, nil, "udp")
	if err != nil {
		utils.LOG.Print("can not connect udp server", dest, err)
		return
	}
	defer con.Close()

	_, err = con.Write(o)
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

	out, err := I.Compress(buf[:n], security.MakeCompressKey(pw))
	if err != nil {
		utils.LOG.Print("compress failed", err)
		return
	}
	_, err = root.WriteToUDP(out, addr)

	utils.LOG.Print("write to client ", err)

}
