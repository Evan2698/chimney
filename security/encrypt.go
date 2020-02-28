package security

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"

	"github.com/Evan2698/chimney/utils"
)

const (
	CHACHA_20 = "chacha20" //CHACHA_20
	GCM       = "gcm"      //GCM
	RAW       = "raw"      //RAW
	PLOY1305  = "p1305"    //PLOY1305
)

const (
	CHACHA_INT   = 5141
	GCM_INT      = 1302
	RAW_INT      = 24869
	PLOY1305_INT = 9011
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

	//GetSize
	GetSize() int
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

	switch name {
	case CHACHA_20:
		i = &cha20{
			name: CHACHA_20,
		}
	case GCM:
		i = &gcm{
			name: GCM,
		}
	case RAW:
		i = &rawS{
			name: RAW,
		}
	case PLOY1305:
		i = &ploy{
			name: PLOY1305,
		}
	default:
		i = nil
	}
	if i != nil {
		i.SetIV(i.MakeSalt())
	}
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
	i := NewEncryptyMethod(name)
	c := make([]byte, len(iv))
	copy(c, iv)
	i.SetIV(c)
	return i
}

//FromByte ...
func FromByte(buf []byte) (EncryptThings, error) {
	if buf == nil {
		return nil, errors.New("invalid paramter")
	}

	op := bytes.NewBuffer(buf)

	var name string
	l := op.Next(2)

	if len(l) < 1 {
		return nil, errors.New("out of length")
	}
	flag := utils.Bytes2Uint16(l)
	name = ""
	switch flag {
	case CHACHA_INT:
		name = CHACHA_20
	case GCM_INT:
		name = GCM
	case RAW_INT:
		name = RAW
	case PLOY1305_INT:
		name = PLOY1305
	default:
		return nil, errors.New("out of length")
	}
	lvl := op.Next(1)
	if len(lvl) < 1 {
		return nil, errors.New("out of length")
	}
	iv := []byte{}
	value := int(lvl[0])
	if value > 0 {
		iv = op.Next(value)
	}
	return NewEncryptyMethodWithIV(name, iv), nil
}

//ToBytes ...
func ToBytes(I EncryptThings) []byte {

	var op bytes.Buffer
	if I.GetName() == CHACHA_20 {
		op.Write([]byte{0x14, 0x15})
	} else if I.GetName() == GCM {
		op.Write([]byte{0x05, 0x16})
	} else if I.GetName() == PLOY1305 {
		op.Write([]byte{0x23, 0x33})
	} else {
		op.Write([]byte{0x61, 0x25})
	}
	lv := (byte)(len(I.GetIV()))
	op.WriteByte(lv)
	if lv > 0 {
		op.Write(I.GetIV())
	}
	return op.Bytes()
}
