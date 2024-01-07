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
	assert.Equal(0, mm.Size("unknown"))

	mm.Clear()

	assert.False(mm.Contain("test"))
	assert.Equal([]string{}, mm.Get("test"))

	mm.Add("test", "1", "2", "3")
	mm.Add("test", "3", "4", "5")

	assert.True(mm.Contain("test"))
	assert.Equal([]string{"1", "2", "3", "3", "4", "5"}, mm.Get("test"))
	assert.Equal(6, mm.Size("test"))

	mm.Remove("test")
	mm.Remove("unknown")

	assert.False(mm.Contain("test"))
	assert.Equal([]string{}, mm.Get("test"))

	mm.Add("test1", "1", "2", "3")
	mm.Add("test2", "3", "4", "5")

	assert.Equal(map[string][]string{
		"test1": {"1", "2", "3"},
		"test2": {"3", "4", "5"},
	}, mm.GetAll())
}
