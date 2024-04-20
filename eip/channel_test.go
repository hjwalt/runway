package eip_test

import (
	"log/slog"
	"testing"

	"github.com/hjwalt/runway/eip"
	"github.com/stretchr/testify/assert"
)

func TestChannel(t *testing.T) {
	assert := assert.New(t)

	testChannel, channelErr := eip.ChannelFromUrl[[]byte]("kafka://localhost:9092?consumer_group=test-group&topic=test-topic")
	assert.NoError(channelErr)

	slog.Info("channel param", "param", testChannel.AllParam())

	assert.Equal("test-group", testChannel.GetParam("consumer_group"))
	assert.Equal("test-topic", testChannel.GetParam("topic"))
	assert.Equal("kafka", testChannel.GetScheme())
	assert.Equal("localhost:9092", testChannel.GetHost())
	assert.Equal("localhost", testChannel.GetHostname())
	assert.Equal("9092", testChannel.GetPort())
}
