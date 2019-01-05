package security

import (
	"socks5/security"
	"testing"
)

func Test_KIL(t *testing.T) {

	v := security.NewEncryptyMethod("chacha20")
	b := security.ToBytes(v)

	t.Log(b)

	c, e := security.FromByte(b)

	t.Log(c.GetName(), c.GetIV())
	t.Log(e)
	t.Log(security.ToBytes(c))
}
