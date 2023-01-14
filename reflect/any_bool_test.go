package reflect_test

import (
	"testing"

	"github.com/hjwalt/runway/reflect"
	"github.com/stretchr/testify/assert"
)

type StructForTest struct {
	Message string
}

func TestBoolToBool(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(true, reflect.GetBool(true))
	assert.Equal(false, reflect.GetBool(false))
}

func TestNilToBool(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(false, reflect.GetBool(nil))
}

func TestPtrToBool(t *testing.T) {
	assert := assert.New(t)

	input := true
	assert.Equal(true, reflect.GetBool(&input))
}

func TestIntToBool(t *testing.T) {
	assert := assert.New(t)

	input := int(1)
	assert.Equal(true, reflect.GetBool(input))
}

func TestIntZeroToBool(t *testing.T) {
	assert := assert.New(t)

	input := int(0)
	assert.Equal(false, reflect.GetBool(input))
}

func TestIntOthersToBool(t *testing.T) {
	assert := assert.New(t)

	input := int(1234)
	assert.Equal(false, reflect.GetBool(input))
}

func TestUintToBool(t *testing.T) {
	assert := assert.New(t)

	input := uint(1)
	assert.Equal(true, reflect.GetBool(input))
}

func TestUintZeroToBool(t *testing.T) {
	assert := assert.New(t)

	input := uint(0)
	assert.Equal(false, reflect.GetBool(input))
}

func TestUintOthersToBool(t *testing.T) {
	assert := assert.New(t)

	input := uint(1234)
	assert.Equal(false, reflect.GetBool(input))
}

func TestStringToBool(t *testing.T) {
	assert := assert.New(t)

	input := "tRuE"
	assert.Equal(true, reflect.GetBool(input))
}

func TestStringFalseToBool(t *testing.T) {
	assert := assert.New(t)

	input := "false"
	assert.Equal(false, reflect.GetBool(input))
}

func TestStringUnknownToBool(t *testing.T) {
	assert := assert.New(t)

	input := "random"
	assert.Equal(false, reflect.GetBool(input))
}

func TestStructToBool(t *testing.T) {
	assert := assert.New(t)

	input := StructForTest{Message: "test"}
	assert.Equal(false, reflect.GetBool(input))
}
