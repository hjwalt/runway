package format_test

import (
	"testing"

	"github.com/hjwalt/runway/format"
	"github.com/stretchr/testify/assert"
)

func TestStringFormat(t *testing.T) {
	assert := assert.New(t)

	f := format.String()

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
