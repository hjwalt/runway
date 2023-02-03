package reflect_test

import (
	"testing"

	"github.com/hjwalt/runway/reflect"
	"github.com/stretchr/testify/assert"
)

func TestUintToUint(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(uint(0), reflect.GetUint(uint(0)))
	assert.Equal(uint8(8), reflect.GetUint8(uint(8)))
	assert.Equal(uint16(16), reflect.GetUint16(uint(16)))
	assert.Equal(uint32(32), reflect.GetUint32(uint(32)))
	assert.Equal(uint64(64), reflect.GetUint64(uint(64)))

	assert.Equal(uint64(0), reflect.GetUint64(uint(0)))
	assert.Equal(uint64(8), reflect.GetUint64(uint8(8)))
	assert.Equal(uint64(16), reflect.GetUint64(uint16(16)))
	assert.Equal(uint64(32), reflect.GetUint64(uint32(32)))
	assert.Equal(uint64(64), reflect.GetUint64(uint64(64)))
}

func TestIntToUint(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(uint64(0), reflect.GetUint64(int(0)))
	assert.Equal(uint64(8), reflect.GetUint64(int8(8)))
	assert.Equal(uint64(16), reflect.GetUint64(int16(16)))
	assert.Equal(uint64(32), reflect.GetUint64(int32(32)))
	assert.Equal(uint64(64), reflect.GetUint64(int64(64)))
}

func TestFloatToUint(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(uint64(32), reflect.GetUint64(float32(32)))
	assert.Equal(uint64(64), reflect.GetUint64(float64(64)))
}

func TestNilToUint(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(uint64(0), reflect.GetUint64(nil))
}

func TestPtrToUint(t *testing.T) {
	assert := assert.New(t)

	input := uint(100)
	var inputNil *uint

	assert.Equal(uint64(100), reflect.GetUint64(&input))
	assert.Equal(uint64(0), reflect.GetUint64(inputNil))
}

func TestBoolToUint(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(uint64(0), reflect.GetUint64(false))
	assert.Equal(uint64(1), reflect.GetUint64(true))
}

func TestStringToUint(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(uint64(0), reflect.GetUint64(""))
	assert.Equal(uint64(123), reflect.GetUint64("123"))
	assert.Equal(uint64(0), reflect.GetUint64("some random value"))
}

func TestStructToUint(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(uint64(0), reflect.GetUint64(StructForTest{Message: "123"}))
}

func TestByteToUInt(t *testing.T) {
	assert := assert.New(t)

	valueBytes := make([]byte, 8)
	reflect.Endian().PutUint64(valueBytes, 8)

	assert.Equal(uint64(8), reflect.GetUint64(valueBytes))
}
