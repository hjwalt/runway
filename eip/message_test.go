package eip_test

import (
	"testing"

	"github.com/hjwalt/runway/eip"
	"github.com/stretchr/testify/assert"
)

func TestMessage(t *testing.T) {
	assert := assert.New(t)

	m := eip.NewMessage[int]()

	m.SetHeader(eip.HeaderId, "1")
	m.SetHeader(eip.HeaderKey, "2")
	m.AddHeader(eip.HeaderId, "3")
	m.AddHeader("nonsense", "4")

	assert.Equal("1", m.GetHeader(eip.HeaderId))
	assert.Equal("2", m.GetHeader(eip.HeaderKey))
	assert.Equal([]string{"1", "3"}, m.ValHeader(eip.HeaderId))

	newheaders := eip.MessageHeader{}
	newheaders.Add("map1", "3")
	newheaders.Add("map2", "4")
	m.PutHeader(newheaders)

	assert.Equal("3", m.GetHeader("map1"))
	assert.Equal("4", m.GetHeader("map2"))
	assert.Equal("", m.GetHeader("1"))

	m.DelHeader("map1")
	assert.Equal("", m.GetHeader("map1"))

	currMap := eip.MessageHeader{}
	currMap.Add("map2", "4")
	assert.Equal(currMap, m.AllHeader())

	assert.Equal(0, m.GetBody())
	m.SetBody(1)
	assert.Equal(1, m.GetBody())
	m.SetBody(2)
	assert.Equal(2, m.GetBody())
}
