package core

import (
	"net"
	"testing"
)

func TestHello12(t *testing.T) {

	remote, err := net.Dial("tcp", "www.baidu.com:443")
	if err != nil {
		t.Log(err)

	}
	a := "sahhghs\n\rzhan"
	t.Log([]byte(a[:]))

	t.Log("remote is:", remote)
	remote.Close()
}

func TestHello45(t *testing.T) {

	ip := net.ParseIP("10.10.10.10").To4()
	t.Log(ip)
}
