package optional_test

import (
	"testing"

	"github.com/hjwalt/runway/optional"
	"github.com/stretchr/testify/assert"
)

func TestOptional(t *testing.T) {
	assert := assert.New(t)

	normalOptional := optional.Of("test")

	assert.True(normalOptional.IsPresent())
	assert.Equal("test", normalOptional.Get())
	assert.Equal("test", normalOptional.GetOrDefault("default"))

	pointerOptional := optional.Empty[string]()
	assert.False(pointerOptional.IsPresent())
	assert.Equal("", pointerOptional.Get())
	assert.Equal("default", pointerOptional.GetOrDefault("default"))

	var strptr *string

	pointerOptional = optional.Pointer(strptr)
	assert.False(pointerOptional.IsPresent())
	assert.Equal("", pointerOptional.Get())
	assert.Equal("default", pointerOptional.GetOrDefault("default"))

	strval := "newval"
	strptr = &strval

	// should not affect the optional
	assert.False(pointerOptional.IsPresent())
	assert.Equal("", pointerOptional.Get())
	assert.Equal("default", pointerOptional.GetOrDefault("default"))

	pointerOptional = optional.Pointer(strptr)
	assert.True(pointerOptional.IsPresent())
	assert.Equal("newval", pointerOptional.Get())
	assert.Equal("newval", pointerOptional.GetOrDefault("default"))
}
