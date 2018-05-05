package sercurity

import (
	"errors"
	"crypto/sha256"
	"io"
	"crypto/rand"
	"crypto/cipher"
	"crypto/aes"
	"github.com/Evan2698/chimney/utils"
	"crypto/sha1"
	"encoding/hex"
	"strings"
	"crypto/hmac"
	"github.com/Yawning/chacha20"	
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


func CompressWithChaCha20(src []byte, iv []byte, key []byte) ([] byte, error) {
   if len(iv) != 8 || len(key) !=32 || len(src) == 0 {
	   return nil, errors.New("parameter is invalid.")
   }

   dst := make([]byte, len(src))

   a, err := chacha20.NewCipher(key, iv)
   if err != nil {
	   return nil, err
   }

   a.XORKeyStream(dst, src)

   return dst, nil
}


func DecompressWithChaCha20(src []byte, iv []byte, key []byte) ([] byte, error) {
	
	return CompressWithChaCha20(src, iv, key)
}
