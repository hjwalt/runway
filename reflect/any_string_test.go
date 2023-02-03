package reflect_test

import (
	"testing"

	"github.com/hjwalt/runway/reflect"
	"github.com/stretchr/testify/assert"
)

func TestStringToString(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("test", reflect.GetString("test"))
}

func TestNilToString(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("", reflect.GetString(nil))
}

func TestPtrToString(t *testing.T) {
	assert := assert.New(t)

	inputForPtr := "testptr"
	var inputNil *string

	assert.Equal("testptr", reflect.GetString(&inputForPtr))
	assert.Equal("", reflect.GetString(inputNil))
}

func TestIntToString(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("0", reflect.GetString(int(0)))
	assert.Equal("8", reflect.GetString(int8(8)))
	assert.Equal("16", reflect.GetString(int16(16)))
	assert.Equal("32", reflect.GetString(int32(32)))
	assert.Equal("64", reflect.GetString(int64(64)))
}

func TestUintToString(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("0", reflect.GetString(uint(0)))
	assert.Equal("8", reflect.GetString(uint8(8)))
	assert.Equal("16", reflect.GetString(uint16(16)))
	assert.Equal("32", reflect.GetString(uint32(32)))
	assert.Equal("64", reflect.GetString(uint64(64)))
}

func TestFloatToString(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("32", reflect.GetString(float32(32)))
	assert.Equal("64", reflect.GetString(float64(64)))
}

func TestBoolToString(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("false", reflect.GetString(false))
	assert.Equal("true", reflect.GetString(true))
}

func TestStructToString(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("", reflect.GetString(StructForTest{Message: "123"}))
}
