package structure_test

import (
	"testing"

	"github.com/hjwalt/runway/structure"
	"github.com/stretchr/testify/assert"
)

func TestMultiMap(t *testing.T) {
	assert := assert.New(t)

	mm := structure.NewMultiMap[string, string]()

	mm.Add("test", "1", "2", "3")

	assert.True(mm.Contain("test"))
	assert.Equal([]string{"1", "2", "3"}, mm.Get("test"))

	assert.False(mm.Contain("unknown"))
	assert.Equal([]string{}, mm.Get("unknown"))

	mm.Clear()

	assert.False(mm.Contain("test"))
	assert.Equal([]string{}, mm.Get("test"))

	mm.Add("test", "1", "2", "3")
	mm.Add("test", "3", "4", "5")

	assert.True(mm.Contain("test"))
	assert.Equal([]string{"1", "2", "3", "3", "4", "5"}, mm.Get("test"))

	mm.Remove("test")

	assert.False(mm.Contain("test"))
	assert.Equal([]string{}, mm.Get("test"))

}
