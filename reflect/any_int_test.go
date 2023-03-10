package reflect_test

import (
	"testing"

	"github.com/hjwalt/runway/reflect"
	"github.com/stretchr/testify/assert"
)

func TestIntToInt(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(int(0), reflect.GetInt(int(0)))
	assert.Equal(int8(8), reflect.GetInt8(int(8)))
	assert.Equal(int16(16), reflect.GetInt16(int(16)))
	assert.Equal(int32(32), reflect.GetInt32(int(32)))
	assert.Equal(int64(64), reflect.GetInt64(int(64)))

	assert.Equal(int64(0), reflect.GetInt64(int(0)))
	assert.Equal(int64(8), reflect.GetInt64(int8(8)))
	assert.Equal(int64(16), reflect.GetInt64(int16(16)))
	assert.Equal(int64(32), reflect.GetInt64(int32(32)))
	assert.Equal(int64(64), reflect.GetInt64(int64(64)))
}

func TestUintToInt(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(int64(0), reflect.GetInt64(uint(0)))
	assert.Equal(int64(8), reflect.GetInt64(uint8(8)))
	assert.Equal(int64(16), reflect.GetInt64(uint16(16)))
	assert.Equal(int64(32), reflect.GetInt64(uint32(32)))
	assert.Equal(int64(64), reflect.GetInt64(uint64(64)))
}

func TestNilToInt(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(int64(0), reflect.GetInt64(nil))
}

func TestPtrToInt(t *testing.T) {
	assert := assert.New(t)

	input := int(1)
	var inputNil *int

	assert.Equal(int64(1), reflect.GetInt64(&input))
	assert.Equal(int64(0), reflect.GetInt64(inputNil))
}

func TestFloatToInt(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(int64(3), reflect.GetInt64(float32(2.6)))
	assert.Equal(int64(3), reflect.GetInt64(float32(2.5)))
	assert.Equal(int64(2), reflect.GetInt64(float64(2.4)))
}

func TestBoolToInt(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(int64(0), reflect.GetInt64(false))
	assert.Equal(int64(1), reflect.GetInt64(true))
}

func TestStringToInt(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(int64(0), reflect.GetInt64(""))
	assert.Equal(int64(123), reflect.GetInt64("123"))
	assert.Equal(int64(0), reflect.GetInt64("some random value"))
}

func TestStructToInt(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(int64(0), reflect.GetInt64(StructForTest{Message: "123"}))
}

func TestByteToInt(t *testing.T) {
	assert := assert.New(t)

	valueBytes := make([]byte, 8)
	reflect.Endian().PutUint64(valueBytes, 8)

	assert.Equal(int64(8), reflect.GetInt64(valueBytes))
	assert.Equal(int64(0), reflect.GetInt64([]byte{}))
}
