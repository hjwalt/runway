package format

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"

	"github.com/hjwalt/runway/logger"
	"github.com/hjwalt/runway/trusted"
)

type AesGcmFormat struct {
	cipher cipher.AEAD
}

// Encrypt
func (helper AesGcmFormat) Mask(value []byte) ([]byte, error) {
	nonce := make([]byte, helper.cipher.NonceSize())

	_, err := io.ReadFull(rand.Reader, nonce)
	logger.ErrorIfErr("error with nonce generation", err)

	outputBytes := helper.cipher.Seal(nonce, nonce, value, nil)
	return outputBytes, nil
}

// Decrypt
func (helper AesGcmFormat) Unmask(value []byte) ([]byte, error) {

	nonceSize := helper.cipher.NonceSize()
	if len(value) < nonceSize {
		return []byte{}, ErrAesGcmDecryptTooShort
	}

	nonce, ciphertext := value[:nonceSize], value[nonceSize:]
	plainBytes, err := helper.cipher.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return []byte{}, errors.Join(ErrAesGcmDecrypt, err)
	}

	return plainBytes, nil
}

func AesGcm(key string) (Mask, error) {
	aescipher, err := aes.NewCipher([]byte(key))
	if err != nil {
		return Plain(), errors.Join(ErrAesCreate, err)
	}

	gcmpad, err := cipher.NewGCM(aescipher)
	gcmpad = trusted.Must(gcmpad, err)

	return AesGcmFormat{cipher: gcmpad}, nil
}

func AesGcmByteKey(key []byte) (Mask, error) {
	aescipher, err := aes.NewCipher(key)
	if err != nil {
		return Plain(), errors.Join(ErrAesCreate, err)
	}

	gcmpad, err := cipher.NewGCM(aescipher)
	gcmpad = trusted.Must(gcmpad, err)

	return AesGcmFormat{cipher: gcmpad}, nil
}

var (
	ErrAesCreate             = errors.New("AES cipher creation error")
	ErrGcmCreate             = errors.New("GCM cipher mode creation error")
	ErrAesGcmEncrypt         = errors.New("AES GCM encrypt error")
	ErrAesGcmDecrypt         = errors.New("AES GCM decrypt error")
	ErrAesGcmDecryptTooShort = errors.New("AES GCM bytes to decrypt too short")
)
