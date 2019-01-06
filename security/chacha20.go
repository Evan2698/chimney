package security

import (
	"crypto/rand"
	"errors"
	"io"
	"time"

	"github.com/Evan2698/chimney/utils"
	"github.com/Yawning/chacha20"
)

type cha20 struct {
	name string
	iv   []byte
}

func (chacha *cha20) Compress(src []byte, key []byte) ([]byte, error) {
	t1 := time.Now() // get current time
	defer func() {
		elapsed := time.Since(t1)
		utils.LOG.Println("takes time(Compress): ", elapsed.String())
	}()

	if len(key) != 32 || len(src) == 0 {
		return nil, errors.New("parameter is invalid")
	}

	dst := make([]byte, len(src))

	a, err := chacha20.NewCipher(key, chacha.iv)
	if err != nil {
		return nil, err
	}

	a.XORKeyStream(dst, src)

	return dst, nil

}

func (chacha *cha20) Uncompress(src []byte, key []byte) ([]byte, error) {
	return chacha.Compress(src, key)
}

func (chacha *cha20) MakeSalt() []byte {
	nonce := make([]byte, 8)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil
	}
	return nonce
}

func (chacha *cha20) GetIV() []byte {
	return chacha.iv
}

func (chacha *cha20) SetIV(iv []byte) {
	chacha.iv = iv
}

func (chacha *cha20) GetName() string {
	return chacha.name
}
