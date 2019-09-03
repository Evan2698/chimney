package core

import (
	"bytes"
	"crypto/hmac"
	"encoding/binary"
	"errors"

	"chimney/security"
	"chimney/utils"
)

// TryparseUDPProtocol ...
func TryparseUDPProtocol(o []byte, pw string) ([]byte, string, security.EncryptThings, error) {
	if len(o) < 32 {
		return nil, "", nil, errors.New("udp data format is incorrect")
	}

	I, err := security.FromByte(o)
	if err != nil {
		utils.LOG.Println("parse security failed", err)
		return nil, "", nil, err
	}

	Ilen := I.GetSize()

	rest := o[Ilen:]
	hlen := rest[0]
	hmac1 := rest[1:(hlen + 1)]

	hmac2 := security.BuildMacHash(I.GetIV(), pw)
	if !hmac.Equal(hmac1, hmac2) {
		utils.LOG.Print("user name verify failed!")
		return nil, "", nil, errors.New("udp verify failed")
	}

	ori, err := I.Uncompress(rest[hlen+1:], security.MakeCompressKey(pw))
	if err != nil {
		utils.LOG.Print("uncompress udp failed!")
		return nil, "", nil, errors.New("uncompress udp failed")
	}

	aLen := binary.BigEndian.Uint32(ori[:4])
	addr := string(ori[4 : aLen+4])
	udpData := ori[4+aLen:]

	return udpData, addr, I, nil
}

// PackUDPData ...
func PackUDPData(pw string, raw []byte) ([]byte, security.EncryptThings, error) {

	var buf bytes.Buffer

	I := security.NewEncryptyMethod("chacha20")
	buf.Write(security.ToBytes(I))

	hmac := security.BuildMacHash(I.GetIV(), pw)
	b := byte(len(hmac))
	buf.WriteByte(b)
	buf.Write(hmac)

	en, err := I.Compress(raw, security.MakeCompressKey(pw))
	if err != nil {
		return nil, nil, errors.New("compress udp data failed")
	}
	buf.Write(en)
	return buf.Bytes(), I, nil
}
