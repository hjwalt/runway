package format_test

import (
	"testing"

	"github.com/hjwalt/runway/format"
	"github.com/stretchr/testify/assert"
)

func TestOgerpon(t *testing.T) {
	assert := assert.New(t)

	f := format.Ogerpon()

	v := []byte("test")

	vb, em := f.Marshal(v)
	assert.NoError(em)

	bv, eu := f.Unmarshal(vb)
	assert.NoError(eu)

	assert.Equal(v, bv)
	assert.Equal([]byte{}, f.Default())
}

func TestOgerponError(t *testing.T) {
	assert := assert.New(t)

	f := format.Ogerpon()

	var err error

	_, err = f.Marshal([]byte("wellspring"))
	assert.ErrorIs(err, format.ErrWellspringMask)

	bytesforerr, _ := f.Marshal([]byte("hearthflame"))
	_, err = f.Unmarshal(bytesforerr)
	assert.ErrorIs(err, format.ErrHearthflameMask)
}
