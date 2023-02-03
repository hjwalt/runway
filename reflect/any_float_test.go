package reflect_test

import (
	"testing"

	"github.com/hjwalt/runway/reflect"
	"github.com/stretchr/testify/assert"
)

func TestFloatToFloat(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(float32(0), reflect.GetFloat32(float32(0)))
	assert.Equal(float32(100), reflect.GetFloat32(float32(100)))
	assert.Equal(float64(0), reflect.GetFloat64(float64(0)))
	assert.Equal(float64(100), reflect.GetFloat64(float64(100)))
}

func TestNilToFloat(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(float32(0), reflect.GetFloat32(nil))
	assert.Equal(float64(0), reflect.GetFloat64(nil))
}

func TestPtrToFloat(t *testing.T) {
	assert := assert.New(t)

	input := float32(100)
	var inputNil *float32

	assert.Equal(float32(100), reflect.GetFloat32(&input))
	assert.Equal(float64(100), reflect.GetFloat64(&input))
	assert.Equal(float64(0), reflect.GetFloat64(inputNil))
}

func TestIntToFloat(t *testing.T) {
	assert := assert.New(t)

	input := int32(100)

	assert.Equal(float32(100), reflect.GetFloat32(input))
	assert.Equal(float64(100), reflect.GetFloat64(input))
}

func TestUintToFloat(t *testing.T) {
	assert := assert.New(t)

	input := uint32(100)

	assert.Equal(float32(100), reflect.GetFloat32(input))
	assert.Equal(float64(100), reflect.GetFloat64(input))
}

func TestStringToFloat(t *testing.T) {

	assert := assert.New(t)

	input := "100.123"

	assert.Equal(float32(100.123), reflect.GetFloat32(input))
	assert.Equal(float64(100.123), reflect.GetFloat64(input))
}

func TestStringDefaultToFloat(t *testing.T) {

	assert := assert.New(t)

	input := ""

	assert.Equal(float32(0), reflect.GetFloat32(input))
	assert.Equal(float64(0), reflect.GetFloat64(input))
}

func TestStringUnknownToFloat(t *testing.T) {

	assert := assert.New(t)

	input := "random1234"

	assert.Equal(float32(0), reflect.GetFloat32(input))
	assert.Equal(float64(0), reflect.GetFloat64(input))
}

func TestStructToFloat(t *testing.T) {
	assert := assert.New(t)

	input := StructForTest{Message: "test"}

	assert.Equal(float32(0), reflect.GetFloat32(input))
	assert.Equal(float64(0), reflect.GetFloat64(input))
}
