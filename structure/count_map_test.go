package structure_test

import (
	"testing"

	"github.com/hjwalt/runway/structure"
	"github.com/stretchr/testify/assert"
)

func TestCountMap(t *testing.T) {
	assert := assert.New(t)

	countMap := structure.NewCountMap[string]()

	countMap.Set("test", 100)

	assert.Equal(int64(100), countMap.Get("test"))
	assert.Equal(int64(0), countMap.Get("unknown"))

	countMap.Add("test", 200)
	assert.Equal(int64(300), countMap.Get("test"))

	countMap.Add("test2", 200)
	assert.Equal(int64(200), countMap.Get("test2"))

	countMap.Subtract("test", 400)
	assert.Equal(int64(-100), countMap.Get("test"))

	countMap.Subtract("test3", 400)
	assert.Equal(int64(-400), countMap.Get("test3"))

	countMap.Remove("test")
	countMap.Remove("test3")
	countMap.Remove("unknown")
	assert.Equal(int64(0), countMap.Get("test"))

	countMap.Set("test1", 101)
	countMap.Set("test2", 102)

	assert.Equal(map[string]int64{
		"test1": 101,
		"test2": 102,
	}, countMap.GetAll())

	countMap.Clear()

	assert.False(countMap.Contain("101"))
}
