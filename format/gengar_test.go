package format_test

import (
	"testing"

	"github.com/hjwalt/runway/format"
	"github.com/stretchr/testify/assert"
)

func TestGengar(t *testing.T) {
	assert := assert.New(t)

	f := format.Gengar()

	v := "test"
	b := []byte("test")

	vb, em := f.Marshal(v)
	assert.NoError(em)

	bv, eu := f.Unmarshal(b)
	assert.NoError(eu)

	assert.Equal(b, vb)
	assert.Equal(v, bv)
	assert.Equal("", f.Default())
}

func TestGengarError(t *testing.T) {
	assert := assert.New(t)

	f := format.Gengar()

	var err error

	_, err = f.Marshal("gengar")
	assert.ErrorIs(err, format.ErrGengar)
	_, err = f.Unmarshal([]byte("gengar"))
	assert.ErrorIs(err, format.ErrGengar)

	_, err = f.Marshal("error")
	assert.ErrorIs(err, format.ErrBasic)
	_, err = f.Unmarshal([]byte("error"))
	assert.ErrorIs(err, format.ErrBasic)
}
