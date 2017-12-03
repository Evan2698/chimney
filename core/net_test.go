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