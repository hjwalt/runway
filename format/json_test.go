package format_test

import (
	"testing"

	"github.com/hjwalt/runway/format"
	"github.com/stretchr/testify/assert"
)

func TestJsonFormat(t *testing.T) {
	assert := assert.New(t)

	type TestStruct struct {
		Name  string
		Value int64
	}

	f := format.Json[TestStruct]()
	v := TestStruct{
		Name:  "test",
		Value: 1234567890,
	}
	b := []byte(`{"Name":"test","Value":1234567890}`)

	vb, em := f.Marshal(v)
	assert.NoError(em)

	bv, eu := f.Unmarshal(b)
	assert.NoError(eu)

	assert.Equal(b, vb)
	assert.Equal(v, bv)
}

func TestJsonFormatPointer(t *testing.T) {
	assert := assert.New(t)

	type TestStruct struct {
		Name  string
		Value int64
	}

	f := format.Json[*TestStruct]()
	v := &TestStruct{
		Name:  "test",
		Value: 1234567890,
	}
	b := []byte(`{"Name":"test","Value":1234567890}`)

	vb, em := f.Marshal(v)
	assert.NoError(em)

	bv, eu := f.Unmarshal(b)
	assert.NoError(eu)

	assert.Equal(b, vb)
	assert.Equal(v, bv)
}

func TestJsonFormatEmptyValue(t *testing.T) {
	assert := assert.New(t)

	type TestStruct struct {
		Name  string
		Value int64
	}

	f := format.Json[*TestStruct]()

	vb, em := f.Marshal(nil)
	assert.NoError(em)
	assert.Equal([]byte{}, vb)

	bv, eu := f.Unmarshal(nil)
	assert.NoError(eu)
	assert.Equal(&TestStruct{}, bv)
}
