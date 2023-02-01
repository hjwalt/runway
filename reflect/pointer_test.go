package reflect_test

import (
	"testing"

	"github.com/hjwalt/runway/reflect"
	"github.com/stretchr/testify/assert"
)

func TestIsPointer(t *testing.T) {
	assert := assert.New(t)

	assert.False(reflect.IsPointer("test"))

	pointerVal := "test"
	assert.True(reflect.IsPointer(&pointerVal))

	var testNil *string
	assert.True(reflect.IsPointer(testNil))
}

func TestIsNil(t *testing.T) {
	assert := assert.New(t)

	assert.False(reflect.IsNil("test"))

	pointerVal := "test"
	assert.False(reflect.IsNil(&pointerVal))

	var testNil *string
	assert.True(reflect.IsNil(testNil))
}
