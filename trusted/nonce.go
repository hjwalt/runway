package trusted

import (
	"crypto/rand"
	"io"
)

func Nonce(size int) []byte {
	nonce := make([]byte, size)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err)
	}
	return nonce
}
