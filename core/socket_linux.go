package core

import (
	"context"
	"crypto/tls"
	"errors"
	"io"
	"net"
	"os"
	"strconv"
	"syscall"

	"chimney/config"

	"chimney/utils"

	quic "github.com/lucas-clemente/quic-go"
)

func createclientsocket(p SocketService, network string, app *config.AppConfig) (io.ReadWriteCloser, error) {
	if app.UseQuic {
		socketHost := net.JoinHostPort(app.Server, strconv.Itoa(int(app.QuicPort)))
		tlsConf := &tls.Config{
			InsecureSkipVerify: true,
			NextProtos:         []string{ProtocolName},
		}

		session, err := quic.DialAddr(socketHost, tlsConf, nil)
		if err != nil {
			utils.LOG.Print("create quick socket session failed", err)
			return nil, err
		}
		stream, err := session.OpenStreamSync(context.Background())
		if err != nil {
			session.Close()
			utils.LOG.Print("create quick socket stream failed", err)
			return nil, err
		}
		utils.LOG.Print("create socket(quic) socket success!")
		out := makeQuicSocket(session, stream).setQuickTimeout(app.Timeout)
		utils.LOG.Print("create socket(quic) socket success!2")

		return out, nil
	}

	host := net.JoinHostPort(app.Server, strconv.Itoa(int(app.ServerPort)))
	outcon, err := CreateCommonSocket(host, network, app.Timeout, p)
	utils.LOG.Print("create as common socket!!")
	return outcon, err
}

// CreateCommonSocket ...
func CreateCommonSocket(host string, network string, timeout int, p SocketService) (net.Conn, error) {

	var outcon net.Conn
	var err error
	var hostip net.IP
	var port int
	init := false
	if p != nil {

		if network == "tcp" {
			tcpAddr, err := net.ResolveTCPAddr("tcp", host)
			if err != nil {
				utils.LOG.Print("parse tcp address failed: ", err)
				return nil, err
			}
			hostip = tcpAddr.IP
			port = tcpAddr.Port
		} else if network == "udp" {
			tcpAddr, err := net.ResolveUDPAddr("udp", host)
			if err != nil {
				utils.LOG.Print("parse tcp address failed: ", err)
				return nil, err
			}
			hostip = tcpAddr.IP
			port = tcpAddr.Port
		}

		var sa syscall.Sockaddr
		if hostip.To4() == nil {
			ipa := hostip.To16()
			utils.LOG.Println("I am ipv6 ", ipa, port, len(hostip))
			sa = &syscall.SockaddrInet6{
				Port: port,
				Addr: [16]byte{ipa[0], ipa[1], ipa[2], ipa[3],
					ipa[4], ipa[5], ipa[6],
					ipa[7], ipa[8], ipa[9],
					ipa[10], ipa[11], ipa[12],
					ipa[13], ipa[14], ipa[15]},
			}

		} else {
			ipa := hostip.To4()
			utils.LOG.Println("I am ipv4 ", ipa, port)
			sa = &syscall.SockaddrInet4{
				Port: port,
				Addr: [4]byte{ipa[0], ipa[1], ipa[2], ipa[3]},
			}
		}

		var fd int
		if network == "tcp" {
			fd, err = syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, syscall.IPPROTO_TCP)
		} else {
			fd, err = syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, syscall.IPPROTO_UDP)
		}
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

		if network == "tcp" {
			err = syscall.SetsockoptInt(fd, syscall.IPPROTO_IP, syscall.IP_TOS, 128)
			if err != nil {
				utils.LOG.Print("set opt int failed", err)
				return nil, err
			}
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
		outcon, err = net.Dial(network, host)
	}
	utils.SetSocketTimeout(outcon, timeout)
	return outcon, err

}
