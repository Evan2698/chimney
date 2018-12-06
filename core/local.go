package core

import (
	"errors"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/Evan2698/chimney/sercurity"
	"github.com/Evan2698/chimney/utils"
)

func connect_server(remote net.Conn, config *AppConfig) (iv []byte, err error) {

	n, err := remote.Write([]byte{0x05, 0x01})
	if err != nil || n < 0 {
		utils.Logger.Print("can not write 0x5 to server!")
		return nil, errors.New("can not write socks 5 flag to server")
	}

	buf := make([]byte, 256)
	n, err = remote.Read(buf)
	if err != nil || buf[0] != 0x5 || buf[1] != 0x2 {
		utils.Logger.Print("can not from server", err)
		return nil, errors.New("no response for authen")
	}

	//utils.Logger.Print("Server Buf  ", buf)

	iv = buf[2 : 2+12]

	utils.Logger.Print("IV  ", iv)

	enCode, err := sercurity.Compress([]byte(config.Password), iv, sercurity.MakeCompressKey(config.Password))
	if err != nil {
		utils.Logger.Print("can not encrypt password!!!")
		return nil, errors.New("can not encrypt password!!!")
	}

	//utils.Logger.Print("encode  ", enCode, len(enCode))

	hmac := sercurity.MakeMacHash(iv, config.Password)

	//utils.Logger.Print("HashHAMC  ", hmac, len(hmac))

	outLen := 2 + len(enCode) + 1 + len(hmac)

	outBuf := make([]byte, outLen)

	outBuf[0] = 0x5
	outBuf[1] = (byte)(len(enCode))
	copy(outBuf[2:2+len(enCode)], enCode)
	outBuf[2+len(enCode)] = (byte)(len(hmac))
	copy(outBuf[2+len(enCode)+1:], hmac)

	utils.Logger.Print("send content  ", len(outBuf))

	n, err = remote.Write(outBuf)
	if err != nil {
		return nil, errors.New("write password to server ")
	}

	n, err = remote.Read(buf)
	if err != nil {
		return nil, errors.New("can not get authen reponse from  server!.")
	}

	if n > 1 && buf[1] != 0 {
		utils.Logger.Println("password incorrect, can not connect server ++++++++!!!!!!", buf[1])
		os.Exit(1)
	}
	utils.Logger.Println("client validation!+++++====!!!^.^")

	return iv, nil
}

func handle_local_server(someone net.Conn, config *AppConfig, iv []byte, remote net.Conn) (*ConnectInfo, error) {

	info := &ConnectInfo{
		cmd: CMD_CONNECT,
	}

	buf := make([]byte, 264)
	n, err := someone.Read(buf)
	if err != nil || n < 0 {
		utils.Logger.Print("read from server failed!", err)
		return info, err
	}

	utils.Logger.Print("CLIENT:   ", buf[:n])

	if n < 2 {
		utils.Logger.Print("methed request error!!!")
		return info, errors.New("methed request error")
	}

	info.cmd = (int)(buf[1])

	if info.cmd != CMD_CONNECT && info.cmd != CMD_UDPASSOCIATE {
		utils.Logger.Print("can not support it")
		return info, errors.New("the method server can not support")
	}

	addressLen := n - 2 - 4
	content := buf[4 : 4+addressLen]

	utils.Logger.Print("domain|", content, "|")

	encode, err := sercurity.Compress(content, iv, sercurity.MakeCompressKey(config.Password))
	if err != nil {
		someone.Write([]byte{0x05, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
		utils.Logger.Print("encrypt the content failed!!!", err)
		return info, errors.New("encrypt the content failed")
	}

	out := make([]byte, len(encode)+2+4+1)
	copy(out[0:4], buf[0:4])
	out[4] = byte(len(encode))
	copy(out[5:5+len(encode)], encode)
	copy(out[5+len(encode):], buf[n-2:])

	utils.Logger.Print("new____+++", out)

	n, err = remote.Write(out)
	if err != nil {
		someone.Write([]byte{0x05, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
		utils.Logger.Print("write domain failed!!!", err)
		return info, errors.New("write domain failed")
	}

	n, err = remote.Read(buf)
	if err != nil || n < 0 {
		someone.Write([]byte{0x05, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
		utils.Logger.Print("server connect response failed!!!", err)
		return info, errors.New("server conect result failed")
	}

	if buf[1] != 0x00 {
		someone.Write([]byte{0x05, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
		utils.Logger.Print("can not support it!!!")
		return info, errors.New("server connect failed, but response return back")
	}

	if info.cmd != CMD_UDPASSOCIATE {
		// TCP CONECTION
		b := make([]byte, 10)
		copy(b[0:8], []byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00})
		b[8] = byte(config.LocalPort & 0xff)
		b[9] = byte((config.LocalPort >> 8) & 0xff)
		n, err = someone.Write(b)
		if err != nil || n < 0 {
			utils.Logger.Print("write to client response failed", err)
			return info, errors.New("write to client response failed")
		}
	} else {
		// UDP associate SUCCESS
		if GPortQueue.IsEmpty() {
			someone.Write([]byte{0x05, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
			utils.Logger.Print("can not get port for this UDP server!!!")
			return nil, errors.New("can not get port for this UDP server")
		}

		udpPort := (uint16)(GPortQueue.Dequeue())
		info.udpport = udpPort
		utils.Logger.Print("alloc port for UDP: ", udpPort)

		ip := net.ParseIP(config.LocalAddress).To4()
		info.udpConnect, err = createUDPListen(ip.String(), info.udpport)
		if err != nil {
			GPortQueue.Enqueue((Item)(info.udpport))
			utils.Logger.Print("create UDP listen failed ", ip.String(), ":", info.udpport)
			return nil, errors.New("create UDP listen failed " + ip.String() + ":" + strconv.Itoa((int)(info.udpport)))
		}

		info.addr = make([]byte, len(ip))
		copy(info.addr, ip)
		info.atype = IsIP4(ip)
		full := append([]byte{0x05, 0x00, 0x00}, info.atype)
		full = append(full, ip...)
		ports := utils.Port2Bytes(udpPort)
		full = append(full, ports...)
		utils.Logger.Print("response the CLIENT UDP SUCCESS", full)
		n, err = someone.Write(full)
		if err != nil || n < 0 {
			utils.Logger.Print("udp packet response failed", err)
			info.udpConnect.Close()
			GPortQueue.Enqueue((Item)(info.udpport))
			return nil, errors.New("write to client response failed")
		}
	}
	return info, nil
}

func hand_local_routine(someone net.Conn, config *AppConfig) {
	utils.Logger.Print("route begin|||||||||||||||||||||||||||||||||||||")
	defer func() {
		someone.Close()
		utils.Logger.Print("raw socket is closed!!!!!")
		utils.Logger.Print("route end!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
	}()

	t1 := time.Now()
	utils.SetReadTimeOut(someone, config.Timeout)

	buf := make([]byte, 256)
	n, err := someone.Read(buf)
	if err != nil || 0 >= n || 0x5 != buf[0] {
		utils.Logger.Print("read error form", err, "read bytes:", n)
		return
	}

	n, err = someone.Write([]byte{0x05, 0x00})
	if err != nil {
		utils.Logger.Print("write to brower or other failed!", err)
		return
	}

	conSocket, err := Build_low_socket(config.Server, config.ServerPort)
	if err != nil {
		utils.Logger.Print("can not connect server", err)
		return
	}

	remote := conSocket.Remote
	defer func() {
		conSocket.Close()
		utils.Logger.Print("remote socket is closed!!!!!!!")
	}()

	utils.SetReadTimeOut(remote, config.Timeout)

	iv, err := connect_server(remote, config)
	if err != nil {
		utils.Logger.Print("can not connect server", err)
		return
	}

	info, err := handle_local_server(someone, config, iv, remote)
	if err != nil {
		utils.Logger.Print("can not handle brower and server!!", err)
		return
	}

	if info.cmd == CMD_CONNECT {

		ssl := NewSSocket(remote, config.Password, iv)

		go Copy_RAW2C(ssl, someone, nil)

		Copy_C2RAW(ssl, someone, nil)
	} else if info.cmd == CMD_UDPASSOCIATE {
		ch := make(chan string)
		udpsocket := NewUDPSocket(remote, config, iv, info, ch)
		err := udpsend(udpsocket)
		utils.Logger.Print("udp error: ", err)
	}

	elapsed := time.Since(t1)
	utils.Logger.Print("takes time:---------------", elapsed)

}

func Run_Local_routine(config *AppConfig) {

	all, err := net.Listen("tcp", config.LocalAddress+":"+strconv.Itoa(config.LocalPort))
	if err != nil {
		utils.Logger.Print("local listen on   ip:port 127.0.0.1: failed!", config.LocalAddress, strconv.Itoa(config.ServerPort), err)
		os.Exit(1)
	}

	defer all.Close()

	for {
		someone, err := all.Accept()
		if err != nil {
			utils.Logger.Print("remote socket failed to open", err)
			continue
		}

		go hand_local_routine(someone, config)
	}

}

// GPortQueue ..
var GPortQueue ItemQueue

// PORTBASE ...
const PORTBASE uint16 = 49800

//InitializePort ..
func InitializePort() int {
	GPortQueue.New()
	var i uint16
	for i = 0; i < 200; i = i + 1 {
		GPortQueue.Enqueue((Item)(PORTBASE + i))
	}

	return 0

}
