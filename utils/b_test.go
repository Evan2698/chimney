package utils

import(
	"testing"
)

func TestHello(t *testing.T){

	var a uint32
	a = 0x12345678
	b := Int2byte(a)
	c := Byte2int(b)

	t.Log("hello", a)
	t.Log(b)
	t.Log("CCC2", c)
}