package core

import (
	"climbwall/utils"
	"net"
	"os"
	"syscall"
)

func create_low_socket(ipString string, port int) (net.Conn, int, error) {

	init := false
	ip := net.ParseIP(ipString)
	ipa := ip.To4()

	sa := &syscall.SockaddrInet4{
		Port: port,
		Addr: [4]byte{ipa[0], ipa[1], ipa[2], ipa[3]},
	}

	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, syscall.IPPROTO_TCP)
	if err != nil {
		utils.Logger.Println("create socket Fd failed", ipString)
		return nil, nil, err
	}

	defer func() {
		if (!init) {
			syscall.Close(fd)
		}
	}

	err = protect_socket(fd)
	if err != nil {
		utils.Logger.Println("protect fd failed:  fd=", fd)
		return nil, nil, err
	}

	err = syscall.Connect(fd, sa)
	if err != nil {
		utils.Logger.Println("connect the ip failed[ " + ipString + " ]")
		return nil, nil, err
	}

	file := os.NewFile(uintptr(fd), "")
	defer file.Close()

	conn, err := net.FileConn(file)
	if err != nil {
		utils.Logger.Println("Create File object failed[ " + ipString + " ]")
		return nil, nil, err
	}

	init = true
	return conn, fd, nil
}

func protect_socket(fd int) error {

	return nil
}
