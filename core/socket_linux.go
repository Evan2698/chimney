package core

import (
	"errors"
	"net"
	"os"
	"syscall"

	"github.com/Evan2698/chimney/utils"
)

func createclientsocket(host string, p SocketService) (net.Conn, error) {
	var outcon net.Conn
	var err error
	init := false
	if p != nil {
		tcpAddr, err := net.ResolveTCPAddr("tcp", host)
		if err != nil {
			utils.LOG.Print("parse tcp address failed: ", err)
			return nil, err
		}

		var sa syscall.Sockaddr
		if tcpAddr.IP.To4() == nil {
			ipa := tcpAddr.IP.To16()
			utils.LOG.Println("I am ipv6 ", ipa, tcpAddr.Port, len(tcpAddr.IP))
			sa = &syscall.SockaddrInet6{
				Port: tcpAddr.Port,
				Addr: [16]byte{ipa[0], ipa[1], ipa[2], ipa[3],
					ipa[4], ipa[5], ipa[6],
					ipa[7], ipa[8], ipa[9],
					ipa[10], ipa[11], ipa[12],
					ipa[13], ipa[14], ipa[15]},
			}

		} else {
			ipa := tcpAddr.IP.To4()
			utils.LOG.Println("I am ipv4 ", ipa, tcpAddr.Port)
			sa = &syscall.SockaddrInet4{
				Port: tcpAddr.Port,
				Addr: [4]byte{ipa[0], ipa[1], ipa[2], ipa[3]},
			}
		}

		fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, syscall.IPPROTO_TCP)
		if err != nil {
			utils.LOG.Print("create socket failed", err)
			return nil, err
		}
		defer func() {
			if !init {
				utils.LOG.Println("socket operator failed. so will close it.")
				syscall.Close(fd)
			}
		}()

		v := p.Protect(fd) // protect first!!!
		if !v {
			utils.LOG.Print("protect socket failed: ", v)
			return nil, errors.New("protect socket failed")
		}

		err = syscall.SetsockoptInt(fd, syscall.IPPROTO_IP, syscall.IP_TOS, 128)
		if err != nil {
			utils.LOG.Print("set opt int failed", err)
			return nil, err
		}

		err = syscall.Connect(fd, sa)
		if err != nil {
			utils.LOG.Print("connect failed:", err)
			return nil, err
		}

		file := os.NewFile(uintptr(fd), "")
		outcon, err = net.FileConn(file)
		if err != nil {
			file.Close()
			utils.LOG.Print("convert to FileConn failed:", err)
			return nil, err
		}
		init = true

	} else {
		outcon, err = net.Dial("tcp", host)
	}

	return outcon, err
}
