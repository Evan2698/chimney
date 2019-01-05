package utils

import (
	"bytes"
	"encoding/binary"
)

//Int2Bytes ...
func Int2Bytes(n int) []byte {
	u := int32(n)
	var hello bytes.Buffer
	binary.Write(&hello, binary.BigEndian, u)
	return hello.Bytes()
}

//Bytes2Int ...
func Bytes2Int(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)
	var tmp int32
	binary.Read(bytesBuffer, binary.BigEndian, &tmp)
	return int(tmp)
}
