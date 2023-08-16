package structure_test

import (
	"testing"

	"github.com/hjwalt/runway/structure"
	"github.com/stretchr/testify/assert"
)

func TestSet(t *testing.T) {
	assert := assert.New(t)

	set := structure.NewSet[string]()

	set.Add("test")
	set.Add("test1", "test2")

	assert.True(set.Contain("test"))
	assert.True(set.Contain("test1", "test2"))

	set.Remove("test1")
	assert.True(set.Contain("test", "test2"))
	assert.False(set.Contain("test1"))

	set.Clear()
	assert.False(set.Contain("test"))
	assert.False(set.Contain("test2"))
}
