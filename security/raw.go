package security

import (
	"time"

	"github.com/Evan2698/chimney/utils"
)

type rawS struct {
	name string
	iv   []byte
}

func (raw *rawS) Compress(src []byte, key []byte) ([]byte, error) {
	t1 := time.Now() // get current time
	defer func() {
		elapsed := time.Since(t1)
		utils.LOG.Println("takes time(Compress): ", elapsed.String())
	}()

	return src, nil
}

func (raw *rawS) Uncompress(src []byte, key []byte) ([]byte, error) {
	return raw.Compress(src, key)
}

func (raw *rawS) MakeSalt() []byte {
	return []byte{}
}

func (raw *rawS) GetIV() []byte {
	return raw.iv
}

func (raw *rawS) SetIV(iv []byte) {
	raw.iv = iv
}

func (raw *rawS) GetName() string {
	return raw.name
}

func (raw *rawS) GetSize() int {
	return 2 + 1
}
