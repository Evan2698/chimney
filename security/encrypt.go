package security

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
)

// EncryptThings  for encrypt the every things
type EncryptThings interface {
	// encrypt the bytes
	Compress(src []byte, key []byte) ([]byte, error)
	// descrypt the bytes
	Uncompress(src []byte, key []byte) ([]byte, error)

	//iv
	GetIV() []byte

	// salt
	MakeSalt() []byte

	//SetIV
	SetIV([]byte)

	//GetName
	GetName() string
}

//BuildMacHash ..
func BuildMacHash(key []byte, message string) []byte {
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return h.Sum(nil)
}

//NewEncryptyMethod ..
func NewEncryptyMethod(name string) EncryptThings {

	var i EncryptThings
	if "chacha20" == name {
		i = &cha20{
			name: "chacha20",
		}
	} else {
		i = &gcm{
			name: "gcm",
		}
	}
	i.SetIV(i.MakeSalt())
	return i
}

//MakeCompressKey ...
func MakeCompressKey(srcKey string) []byte {
	r := sha1.Sum([]byte(srcKey))
	out := hex.EncodeToString(r[:])
	out = strings.ToUpper(out)
	return ([]byte(out[:]))[:32]
}

// NewEncryptyMethodWithIV ..
func NewEncryptyMethodWithIV(name string, iv []byte) EncryptThings {

	var i EncryptThings
	if "chacha20" == name {
		i = &cha20{
			name: "chacha20",
		}
	} else {
		i = &gcm{
			name: "gcm",
		}
	}
	i.SetIV(iv)
	return i
}

//FromByte ...
func FromByte(buf []byte) (EncryptThings, error) {
	if buf == nil {
		return nil, errors.New("invalid paramter")
	}

	op := bytes.NewBuffer(buf)

	l := op.Next(1)
	if len(l) < 1 {
		return nil, errors.New("out of length")
	}
	value := int(l[0])

	name := string(op.Next(value))
	lvl := op.Next(1)

	if len(lvl) < 1 {
		return nil, errors.New("out of length")
	}
	value = int(lvl[0])
	iv := op.Next(value)

	return NewEncryptyMethodWithIV(name, iv), nil
}

//ToBytes ...
func ToBytes(I EncryptThings) []byte {

	var op bytes.Buffer

	nalen := byte(len(I.GetName()))
	op.WriteByte(nalen)
	op.WriteString(I.GetName())

	lv := (byte)(len(I.GetIV()))
	op.WriteByte(lv)
	op.Write(I.GetIV())

	return op.Bytes()

}
