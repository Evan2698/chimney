//  +build DRCLO

package core

import (
	"climbwall/unixsocket"
	"climbwall/utils"
	"errors"
	"net"
	"os"
	"syscall"
)

func Build_low_socket(ipString string, port int) (net.Conn, int, error) {

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
		return nil, -1, err
	}

	defer func() {
		if !init {
			syscall.Close(fd)
		}
	}()

	err = protect_socket(fd)
	if err != nil {
		utils.Logger.Println("protect fd failed:  fd=", fd)
		return nil, -1, err
	}

	err = syscall.Connect(fd, sa)
	if err != nil {
		utils.Logger.Println("connect the ip failed[ " + ipString + " ]")
		return nil, -1, err
	}

	file := os.NewFile(uintptr(fd), "")
	defer file.Close()

	conn, err := net.FileConn(file)
	if err != nil {
		utils.Logger.Println("Create File object failed[ " + ipString + " ]")
		return nil, -1, err
	}

	init = true
	return conn, fd, nil
}

func protect_socket(fd int) error {

	conn, err := net.Dial("unix", "protect_path")
	if err != nil {
		utils.Logger.Println("can not create unix socket", err)
		return err
	}

	defer conn.Close()

	fdsock, ok := conn.(*net.UnixConn)
	if !ok {
		utils.Logger.Println("can not create unix socket", err)
		return errors.New("can not create unix socket!!!")
	}

	usock := unixsocket.New(fdsock)

	defer func() {
		usock.Close()
	}()

	err = usock.WriteFD(fd)
	if err != nil {
		utils.Logger.Println("can not create unix socket", err)
		return err
	}

	buf := make([]byte, 1)
	n, err := conn.Read(buf)
	if err != null {
		utils.Logger.Println("failed to send result!!!")
		return err
	}

	if n <= 0 || b[0] != 0 {
		utils.Logger.Println("unix send failed!!", []byte(r))
		return errors.New("unix protect socket failed!!")
	}

	return nil
}

