package sercurity

import (
	"crypto/sha256"
	"io"
	"crypto/rand"
	"crypto/cipher"
	"crypto/aes"
	"climbwall/utils"
	"crypto/sha1"
	"encoding/hex"
	"strings"
	"crypto/hmac"
	
)


func Compress(src []byte, iv []byte, key []byte) ([] byte, error) {


	block, err := aes.NewCipher(key)
	if err != nil {
		utils.Logger.Println("key of AES is invalid!")		
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		utils.Logger.Println("key of AES is invalid!")
		return nil, err
	}

	ciphertext := aesgcm.Seal(nil, iv, src, nil)

	return ciphertext, nil
}


func Uncompress(src []byte, iv []byte, key []byte) ([] byte, error){
	block, err := aes.NewCipher(key)
	if err != nil {
		utils.Logger.Println("key of AES is invalid!(uncompress)")		
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		utils.Logger.Println("key of AES is invalid!(uncompress)")
		return nil, err
	}

	plaintext, err := aesgcm.Open(nil, iv, src, nil)

	return plaintext, err

}

func MakeCompressKey(srcKey string) []byte {
	r := sha1.Sum([]byte(srcKey))
	out := hex.EncodeToString(r[:])
	out = strings.ToUpper(out)
	return ([]byte(out[:]))[:32]
}

func MakeSalt() [] byte {

	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil
	}
	return nonce
}

func MakeMacHash(key []byte, message string) []byte {
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return h.Sum(nil)
}
