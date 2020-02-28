package security

import (
	"crypto/rand"
	"io"
	"time"

	"github.com/Evan2698/chimney/utils"
	"golang.org/x/crypto/chacha20poly1305"
)

type ploy struct {
	name string
	iv   []byte
}

func (p *ploy) Compress(src []byte, key []byte) ([]byte, error) {
	t1 := time.Now() // get current time
	defer func() {
		elapsed := time.Since(t1)
		utils.LOG.Println("takes time(Compress): ", elapsed.String())
	}()

	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		utils.LOG.Println("key of chacha20Poly1305 is invalid!")
		return nil, err
	}

	ciphertext := aead.Seal(nil, p.iv, src, nil)
	return ciphertext, nil
}

func (p *ploy) Uncompress(src []byte, key []byte) ([]byte, error) {
	t1 := time.Now() // get current time
	defer func() {
		elapsed := time.Since(t1)
		utils.LOG.Println("takes time(Uncompress): ", elapsed.String())
	}()

	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		utils.LOG.Println("key of chacha20Poly1305 is invalid!(uncompress)")
		return nil, err
	}

	plaintext, err := aead.Open(nil, p.iv, src, nil)
	return plaintext, err
}

func (p *ploy) MakeSalt() []byte {
	nonce := make([]byte, 24)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil
	}
	return nonce
}

func (p *ploy) GetIV() []byte {
	return p.iv
}

func (p *ploy) SetIV(iv []byte) {
	p.iv = iv
}

func (p *ploy) GetName() string {
	return p.name
}

func (p *ploy) GetSize() int {
	return 2 + 1 + len(p.iv)
}
