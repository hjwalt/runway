package structure_test

import (
	"testing"

	"github.com/hjwalt/runway/structure"
	"github.com/stretchr/testify/assert"
)

func TestVoidCast(t *testing.T) {
	assert := assert.New(t)

	var voidValue structure.Void

	voidValue = 1
	assert.Equal(voidValue, 1)

	intValue := 2
	voidValue = structure.Void(intValue)
	assert.Equal(voidValue, 2)

	castTest, castSuccessful := voidValue.(int)
	assert.Equal(castTest, 2)
	assert.True(castSuccessful)

	voidValue = "string"
	assert.Equal(voidValue, "string")
}
