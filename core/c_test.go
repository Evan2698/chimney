package core

import (
	"socks5/utils"
	"testing"
)

func TestKOP(t *testing.T) {
	abc := utils.Int2Bytes(4)
	t.Log(abc)

	u := utils.Bytes2Int(abc)
	t.Log(u)

}
