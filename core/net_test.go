package core

import (
	"net"
	"testing"
)

func TestHello12(t *testing.T) {

	remote, err := net.Dial("tcp", "www.google.com:443")
	if err != nil {
		t.Log(err)

	}
	t.Log("remote is:", remote)
	remote.Close()
}