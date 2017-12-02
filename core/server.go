package core

import (
	"bytes"
	"syscall"
	"encoding/binary"
	"errors"
	"crypto/hmac"
	"fmt"
	"climbwall/sercurity"
	"climbwall/utils"
	"os"
	"strconv"
	"net"

)


func handshark(someone net.Conn, config * AppConfig, salt []byte) error {

	buf := make([]byte, 530)
	n, err := someone.Read(buf)
	if err != nil || 0 >= n {
		utils.Logger.Fatal("read error form", err, "read bytes:", n)
		return err
	}

	if 0x5 != buf[0] {
		err = fmt.Errorf("can not support this version %X", buf[0])
		utils.Logger.Fatalln(err)
		someone.Write([]byte{0x05, 0xff})
		return err
	}

	// need user name and password 
	n, err = someone.Write([]byte{0x05, 0x02})
	if (err != nil) {
		utils.Logger.Fatalln("can not write Authen method to client", err)
		return err
	}

	n, err = someone.Write(salt)
	if (err != nil) {
		utils.Logger.Fatalln("can not write salt to client", err)
		return err
	}

	n, err = someone.Read(buf)
	if (err != nil || buf[0] != 0x5 ) {

		utils.Logger.Fatalln("can not recieve message", err)
		return err
	}

	var userNameLen uint32
	var passwordLen uint32 

	userNameLen =  (uint32)(buf[1])
	userNamebytes := buf[2: 2 + userNameLen ]
	passwordLen = (uint32)(buf[2 + userNameLen])
	passwordbytes := buf [ 2  + userNameLen + 1 : 2  + userNameLen + 1 + passwordLen ]


	user, err := sercurity.Uncompress(userNamebytes, salt, sercurity.MakeCompressKey(config.Password))
	hmac1 := sercurity.MakeMacHash(salt, string(user))

	if (hmac.Equal(hmac1, passwordbytes)) {
		
		
		if 0 != bytes.Compare([]byte(config.Password), userNamebytes){
			someone.Write([]byte{0x05, 0xff})
			return errors.New("user password content incorrect!")
		}

		n, err = someone.Write([]byte{0x05, 0x00})
		if (err != nil) {
			return errors.New("password response send to client failed!")
		}

	    return nil
	}
	someone.Write([]byte{0x05, 0xff})

	return errors.New("user password incorrect!")

}

func handleConnect(someone net.Conn, config * AppConfig, salt []byte) (addr string, err error){

	buf := make([]byte, 258)	
	n, err := someone.Read(buf)
	if (err != nil || n <0  || buf[0] != 0x05) {
		utils.Logger.Fatalln("can not read remote address from client", err)
		return "", err 
	}

	if (buf[1] != 0x1) {
		utils.Logger.Fatalln("does not support other method except connect")
		return "", errors.New("does not support other method except connect")
	}

	var cLen int
	cLen = (int)(buf[5])
	content, err := sercurity.Uncompress(buf[6 : 6+ cLen], salt, sercurity.MakeCompressKey(config.Password))
	if (err != nil) {
		return "", errors.New("does not parse the address from CC")
	}


	var dIP string
	switch buf[3] &  0xf {
	case 0x01:
		//	IP V4 address: X'01'
		dIP = net.IP(content).String()
	case 0x03:
		//	DOMAINNAME: X'03'
		dIP = net.IP(content).String()
	case 0x04:
		//	IP V6 address: X'04'
		dIP = net.IP(content).String()
	default:
		return "", errors.New("on default, the address is nil!")
	}

	port := binary.BigEndian.Uint16(buf[n-2 : n])
	host := net.JoinHostPort(dIP, strconv.Itoa(int(port)))
	
	return host, nil
}






func handleRoutine(someone net.Conn, config * AppConfig) {

	salt := sercurity.MakeSalt()

	defer someone.Close()

	err := handshark(someone, config, salt)
	if (err != nil) {
		utils.Logger.Fatal("failed handshark.!!!!", err)
		return
	}

	addr, err := handleConnect(someone, config, salt)
	if (addr == "" || err != nil ){
		utils.Logger.Fatal("parse failed", err)
		someone.Write([]byte{0x05, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
		return
	}

	remote, err := net.Dial("tcp", addr)
	if err != nil {
		if ne, ok := err.(*net.OpError); ok && (ne.Err == syscall.EMFILE || ne.Err == syscall.ENFILE) {
			// log too many open file error
			// EMFILE is process reaches open file limits, ENFILE is system limit
			utils.Logger.Fatal("dial error:", err)
		} else {
			utils.Logger.Fatal("error connecting to:", addr, err)
		}

		someone.Write([]byte{0x05, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
		return
	}

	defer remote.Close()

	someone.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

	ssl := NewSSocket(someone, config.Password, salt)

	go func(sic * SSocketWrapper, client net.Conn) {
		for ; ; {
			neterr := sic.Copy2RaW(client)
			if (neterr != nil) {
				utils.Logger.Println("copy completed or failed!")
				break
			}

		}
	}(ssl,remote)


	for ;;  {
		neterr := ssl.WriteFromRaw(remote)
		if (neterr != nil) {
			utils.Logger.Println("write error or write complete!!!")
			break
		}
	}
}



func Run_server_routine(config * AppConfig){

	all, err := net.Listen("tcp", config.Server + ":"+ strconv.Itoa(config.ServerPort))
	if err != nil {	
		utils.Logger.Fatal("can not build server on ip:port", config.Server + ":" + strconv.Itoa(config.ServerPort))
		os.Exit(1)
	}

	defer all.Close()

	for {
		someone, err := all.Accept()
		if err != nil {
			utils.Logger.Fatal("remote socket failed to open", err)
			continue
		}
		
		go handleRoutine(someone, config)
	}

}