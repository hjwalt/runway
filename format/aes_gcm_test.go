package format_test

import (
	"crypto/rand"
	"testing"

	"github.com/hjwalt/runway/format"
	"github.com/stretchr/testify/assert"
)

func TestAesGcm(t *testing.T) {
	assert := assert.New(t)

	aesFormat, err := format.AesGcmMask("asdfasdfasdfasdfasdfasdfasdfasdf")

	assert.NoError(err)

	strBytes := []byte("test")
	encryptedBytes, encryptionErr := aesFormat.Marshal(strBytes)

	assert.NoError(encryptionErr)
	assert.NotEqual(strBytes, encryptedBytes)

	decryptedBytes, decryptionErr := aesFormat.Unmarshal(encryptedBytes)

	assert.NoError(decryptionErr)
	assert.Equal(strBytes, decryptedBytes)
}

func TestAesGcmAesCreateError(t *testing.T) {
	assert := assert.New(t)

	_, err := format.AesGcm("asdfasdfasdfasdfasdfasdfasdf")

	assert.ErrorIs(err, format.ErrAesCreate)
}

func TestAesGcmAesDecryptError(t *testing.T) {
	assert := assert.New(t)

	aesFormat, err := format.AesGcm("asdfasdfasdfasdfasdfasdfasdfasdf")

	assert.NoError(err)

	_, decryptionErr := aesFormat.Unmarshal([]byte("asdfasdfasdfasdfasdfasdfasdf"))
	assert.ErrorIs(decryptionErr, format.ErrAesGcmDecrypt)
}

func TestAesGcmAesDecryptDataTooShort(t *testing.T) {
	assert := assert.New(t)

	aesFormat, err := format.AesGcm("asdfasdfasdfasdfasdfasdfasdfasdf")

	assert.NoError(err)

	_, decryptionErr := aesFormat.Unmarshal([]byte("a"))
	assert.ErrorIs(decryptionErr, format.ErrAesGcmDecryptTooShort)
}

func TestAesGcmRandomBytes(t *testing.T) {
	assert := assert.New(t)

	key := make([]byte, 32)

	_, randErr := rand.Read(key)
	assert.NoError(randErr)

	aesFormat, err := format.AesGcmMaskByteKey(key)

	assert.NoError(err)

	strBytes := []byte("test")
	encryptedBytes, encryptionErr := aesFormat.Marshal(strBytes)

	assert.NoError(encryptionErr)
	assert.NotEqual(strBytes, encryptedBytes)

	decryptedBytes, decryptionErr := aesFormat.Unmarshal(encryptedBytes)

	assert.NoError(decryptionErr)
	assert.Equal(strBytes, decryptedBytes)
}
