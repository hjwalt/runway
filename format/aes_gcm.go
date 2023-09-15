package format

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"

	"github.com/hjwalt/runway/trusted"
)

type AesGcmFormat struct {
	cipher cipher.AEAD
}

func (helper AesGcmFormat) Default() []byte {
	return []byte{}
}

// Encrypt
func (helper AesGcmFormat) Marshal(value []byte) ([]byte, error) {
	nonce := trusted.Nonce(helper.cipher.NonceSize())
	outputBytes := helper.cipher.Seal(nonce, nonce, value, nil)
	return outputBytes, nil
}

// Decrypt
func (helper AesGcmFormat) Unmarshal(value []byte) ([]byte, error) {

	nonceSize := helper.cipher.NonceSize()
	if len(value) < nonceSize {
		return helper.Default(), ErrAesGcmDecryptTooShort
	}

	nonce, ciphertext := value[:nonceSize], value[nonceSize:]
	plainBytes, err := helper.cipher.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return helper.Default(), errors.Join(ErrAesGcmDecrypt, err)
	}

	return plainBytes, nil
}

func AesGcm(key string) (Format[[]byte], error) {
	aescipher, err := aes.NewCipher([]byte(key))
	if err != nil {
		return Bytes(), errors.Join(ErrAesCreate, err)
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
