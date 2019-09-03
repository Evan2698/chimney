package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"time"

	"chimney/utils"
)

type gcm struct {
	name string
	iv   []byte
}

func (g *gcm) Compress(src []byte, key []byte) ([]byte, error) {
	t1 := time.Now() // get current time
	defer func() {
		elapsed := time.Since(t1)
		utils.LOG.Println("takes time(Compress): ", elapsed.String())
	}()

	block, err := aes.NewCipher(key)
	if err != nil {
		utils.LOG.Println("key of AES is invalid!")
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		utils.LOG.Println("key of AES is invalid!")
		return nil, err
	}

	ciphertext := aesgcm.Seal(nil, g.iv, src, nil)
	return ciphertext, nil
}

func (g *gcm) Uncompress(src []byte, key []byte) ([]byte, error) {
	t1 := time.Now() // get current time
	defer func() {
		elapsed := time.Since(t1)
		utils.LOG.Println("takes time(Uncompress): ", elapsed.String())
	}()

	block, err := aes.NewCipher(key)
	if err != nil {
		utils.LOG.Println("key of AES is invalid!(uncompress)")
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		utils.LOG.Println("key of AES is invalid!(uncompress)")
		return nil, err
	}

	plaintext, err := aesgcm.Open(nil, g.iv, src, nil)
	return plaintext, err
}

func (g *gcm) MakeSalt() []byte {
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil
	}
	return nonce
}

func (g *gcm) GetIV() []byte {
	return g.iv
}

func (g *gcm) SetIV(iv []byte) {
	g.iv = iv
}

func (g *gcm) GetName() string {
	return g.name
}

func (g *gcm) GetSize() int {
	return 2 + 1 + len(g.iv)
}
