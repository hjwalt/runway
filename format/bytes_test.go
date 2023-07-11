package format_test

import (
	"testing"

	"github.com/hjwalt/runway/format"
	"github.com/stretchr/testify/assert"
)

func TestBytes(t *testing.T) {
	assert := assert.New(t)
	f := format.Bytes()

	v := []byte("test")
	b := []byte("test")

	vb, em := f.Marshal(v)
	assert.NoError(em)

	bv, eu := f.Unmarshal(b)
	assert.NoError(eu)

	assert.Equal(b, vb)
	assert.Equal(v, bv)
	assert.Equal([]byte{}, f.Default())
}
