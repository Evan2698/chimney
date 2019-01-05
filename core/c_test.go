package core

import (
	"testing"

	"github.com/Evan2698/chimney/utils"
)

func TestKOP(t *testing.T) {
	abc := utils.Int2Bytes(4)
	t.Log(abc)

	u := utils.Bytes2Int(abc)
	t.Log(u)

}
