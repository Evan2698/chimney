//  +build DRCLO

package core

import (
	"errors"
	"net"
	"os"
	"syscall"

	"github.com/Evan2698/chimney/utils"
)

func Build_low_socket(ipString string, port int) (*CommonSocket, error) {

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
		return nil, err
	}

	defer func() {
		if !init {
			syscall.Close(fd)
		}
	}()

	err = protect_socket(fd)
	if err != nil {
		utils.Logger.Println("protect fd failed:  fd=", fd)
		return nil, err
	}

	err = syscall.Connect(fd, sa)
	if err != nil {
		utils.Logger.Println("connect the ip failed[ " + ipString + " ]")
		return nil, err
	}

	file := os.NewFile(uintptr(fd), "")
	defer file.Close()

	conn, err := net.FileConn(file)
	if err != nil {
		utils.Logger.Println("Create File object failed[ " + ipString + " ]")
		return nil, err
	}

	init = true
	return &CommonSocket{
		Remote: conn,
		Fd:     uintptr(fd),
	}, nil
}

func protect_socket(fd int) error {

	/*var path = GUNIXPATH + "/protect_path"
	conn, err := net.Dial("unix", path)
	if err != nil {
		utils.Logger.Println("can not create unix socket", err, path)
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
	if err != nil {
		utils.Logger.Println("failed to send result!!!")
		return err
	}

	if n <= 0 || buf[0] != 0 {
		utils.Logger.Println("unix send failed!!", buf)
		return errors.New("unix protect socket failed!!")
	}*/

	if GSocketInterface == nil {
		return errors.New("must register interface first")
	}

	GSocketInterface.Protect(fd)

	return nil
}

func (con *CommonSocket) Close() {
	con.Remote.Close()
	if con.Fd != 0 {
		syscall.Close(int(con.Fd))
	}
}
