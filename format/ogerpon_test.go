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

	vb, em := f.Mask(v)
	assert.NoError(em)
	assert.Equal("tset", string(vb))

	bv, eu := f.Unmask(vb)
	assert.NoError(eu)

	assert.Equal(v, bv)
}

func TestOgerponError(t *testing.T) {
	assert := assert.New(t)

	f := format.Ogerpon()

	var err error

	_, err = f.Mask([]byte("wellspring"))
	assert.ErrorIs(err, format.ErrWellspringMask)
	_, err = f.Mask([]byte("error"))
	assert.ErrorIs(err, format.ErrBasic)

	_, err = f.Unmask([]byte("emalfhtraeh"))
	assert.ErrorIs(err, format.ErrHearthflameMask)
	_, err = f.Unmask([]byte("rorre"))
	assert.ErrorIs(err, format.ErrBasic)
}
