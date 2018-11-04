package utils

import (
	"bytes"
	"encoding/binary"
)

// Int2byte function
//
func Int2byte(n uint32) []byte {

	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.BigEndian, n)
	return buf.Bytes()
}

// Byte2int function
//
func Byte2int(src []byte) uint32 {
	bytesBuffer := bytes.NewBuffer(src)
	var tmp uint32
	binary.Read(bytesBuffer, binary.BigEndian, &tmp)
	return tmp
}
