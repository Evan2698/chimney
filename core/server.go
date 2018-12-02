package core

import (
	"bytes"
	"crypto/hmac"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"syscall"
	"time"

	"github.com/Evan2698/chimney/sercurity"
	"github.com/Evan2698/chimney/utils"
)

func handshark(someone net.Conn, config *AppConfig, salt []byte) error {

	buf := make([]byte, 540)
	n, err := someone.Read(buf)
	if err != nil || 0 >= n {
		utils.Logger.Print("read error form", err, "read bytes:", n)
		return err
	}

	if 0x5 != buf[0] {
		err = fmt.Errorf("can not support this version %X", buf[0])
		utils.Logger.Print(err)
		someone.Write([]byte{0x05, 0xff})
		return err
	}

	// need user name and password
	buf[0] = 0x5
	buf[1] = 0x2
	copy(buf[2:2+12], salt)
	out := buf[0 : 2+12]
	n, err = someone.Write(out)
	if err != nil {
		utils.Logger.Print("can not write Authen method to client", err)
		return err
	}

	n, err = someone.Read(buf)
	if err != nil || buf[0] != 0x5 {

		utils.Logger.Print("can not recieve message", err)
		return err
	}

	//utils.Logger.Print("RECC  ", buf)

	var userNameLen uint32
	var passwordLen uint32

	userNameLen = (uint32)(buf[1])
	userNamebytes := buf[2 : 2+userNameLen]
	passwordLen = (uint32)(buf[2+userNameLen])
	passwordbytes := buf[2+userNameLen+1 : 2+userNameLen+1+passwordLen]

	utils.Logger.Print("userNameLen  ", userNameLen)
	utils.Logger.Print("passwordLen  ", passwordLen)
	utils.Logger.Print("passwordbytes  ", passwordbytes)
	user, err := sercurity.Uncompress(userNamebytes, salt, sercurity.MakeCompressKey(config.Password))
	if err != nil {
		utils.Logger.Print("can not uncompress user name", err)
		return err
	}

	hmac1 := sercurity.MakeMacHash(salt, string(user))
	//utils.Logger.Print("hmac1  ", hmac1)

	if hmac.Equal(hmac1, passwordbytes) {

		if 0 != bytes.Compare([]byte(config.Password), user) {
			someone.Write([]byte{0x05, 0xff})
			return errors.New("user password content incorrect!")
		}

		n, err = someone.Write([]byte{0x05, 0x00})
		if err != nil {
			return errors.New("password response send to client failed!")
		}

		return nil
	}
	someone.Write([]byte{0x05, 0xff})

	return errors.New("user password incorrect!")

}

func parsehost(t byte, content []byte) string {
	var dIP string
	it := content[:len(content)-2]
	switch t {
	case 0x01:
		//	IP V4 address: X'01'
		dIP = net.IP(it).String()
	case 0x03:
		//	DOMAINNAME: X'03'
		dIP = string(it[1:])
	case 0x04:
		//	IP V6 address: X'04'
		dIP = net.IP(it).String()
	default:
		dIP = "127.0.0.1"
	}

	port := binary.BigEndian.Uint16(content[len(content)-2:])

	host := net.JoinHostPort(dIP, strconv.Itoa(int(port)))

	return host
}

func handleConnect(someone net.Conn, config *AppConfig, salt []byte) (addr *ConnectInfo, err error) {

	R := &ConnectInfo{}
	buf := make([]byte, 300)
	n, err := someone.Read(buf)
	if err != nil || n < 0 || buf[0] != 0x05 {
		utils.Logger.Print("can not read remote address from client", err)
		return nil, err
	}

	R.cmd = int(buf[1])
	if R.cmd != CMD_CONNECT && R.cmd != CMD_UDPASSOCIATE {
		utils.Logger.Print("does not support other method except connect")
		return nil, errors.New("does not support other method except connect")
	}

	utils.Logger.Println("accesss address ", len(buf), err)
	var cLen int
	cLen = (int)(buf[4])
	content, err := sercurity.Uncompress(buf[5:5+cLen], salt, sercurity.MakeCompressKey(config.Password))
	if err != nil {
		return nil, errors.New("does not parse the address from CC")
	}

	R.addr = make([]byte, len(content))
	copy(R.addr, content)
	R.addresstype = buf[3] & 0xf

	utils.Logger.Println("domain content: ", len(content))

	R.host = parsehost(buf[3], append(content, buf[n-2:n]...))

	R.port = binary.BigEndian.Uint16(buf[n-2 : n])

	utils.Logger.Println("host ", R.host)

	return R, nil
}

func createTCPConnect(host string) (net.Conn, error) {

	remote, err := net.Dial("tcp", host)
	if err != nil {
		if ne, ok := err.(*net.OpError); ok && (ne.Err == syscall.EMFILE || ne.Err == syscall.ENFILE) {
			// log too many open file error
			// EMFILE is process reaches open file limits, ENFILE is system limit
			utils.Logger.Print("dial error: ", err)
		} else {
			utils.Logger.Print("error connecting to:  ", host, err)
		}
		return nil, err
	}

	return remote, nil
}

func createDUPConnect(host string) (net.Conn, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", host)
	if err != nil {
		utils.Logger.Print("ResolveUDPAddr failed:  ", host, err)

		return nil, err
	}

	remote, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		utils.Logger.Print("establish UDP Failed:  ", host, err)

		return nil, err
	}

	return remote, nil
}

func createFakeDUPConnect(host string) (net.Conn, error) {

	return nil, nil
}

func handleRoutine(someone net.Conn, config *AppConfig) {

	t1 := time.Now()
	utils.SetReadTimeOut(someone, config.Timeout)
	salt := sercurity.MakeSalt()

	defer someone.Close()

	err := handshark(someone, config, salt)
	if err != nil {
		utils.Logger.Print("failed handshark.!!!!", err)
		return
	}

	info, err := handleConnect(someone, config, salt)
	if info == nil || err != nil || len(info.host) == 0 {
		utils.Logger.Print("parse failed", err)
		someone.Write([]byte{0x05, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
		return
	}

	utils.Logger.Print("address:   |", info.host+"|")

	var remote net.Conn
	if CMD_CONNECT == info.cmd {
		remote, err = createTCPConnect(info.host)
	} else if CMD_UDPASSOCIATE == info.cmd {
		remote, err = createFakeDUPConnect(info.host)
	}
	if err != nil {
		someone.Write([]byte{0x05, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
		utils.Logger.Print("create connection failed!!!")
		return
	}

	defer func() {
		if remote != nil {
			remote.Close()
		}
	}()

	utils.SetReadTimeOut(remote, config.Timeout)

	someone.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

	if info.cmd == CMD_CONNECT {
		ssl := NewSSocket(someone, config.Password, salt)

		go Copy_C2RAW(ssl, remote, nil)

		Copy_RAW2C(ssl, remote, nil)

	} else if CMD_UDPASSOCIATE == info.cmd {
		//TODO
		ch := make(chan string)
		ss := NewUDPSocket(someone, config.Password, salt, info, ch)
		udpserverRoutine(ss, remote)
	}

	elapsed := time.Since(t1)
	utils.Logger.Print("takes time:---------------", elapsed)
}

func Run_server_routine(config *AppConfig) {

	all, err := net.Listen("tcp", config.Server+":"+strconv.Itoa(config.ServerPort))
	if err != nil {
		utils.Logger.Print("can not build server on ip:port  ", config.Server+":"+strconv.Itoa(config.ServerPort))
		utils.Logger.Print("app will exit!!!!")
		os.Exit(1)
	}

	defer all.Close()

	for {
		someone, err := all.Accept()
		if err != nil {
			utils.Logger.Print("remote socket failed to open", err)
			continue
		}

		go handleRoutine(someone, config)
	}

}
