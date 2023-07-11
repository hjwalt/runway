package format_test

import (
	"testing"

	"github.com/hjwalt/runway/format"
	"github.com/hjwalt/runway/logger"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

func TestProtojsonFormat(t *testing.T) {
	assert := assert.New(t)

	f := format.Protojson[*Message]()
	v := &Message{
		Topic:  "test",
		Offset: 123456789,
	}
	b1 := []byte(`{"topic":"test", "offset":"123456789"}`)
	b2 := []byte(`{"topic":"test","offset":"123456789"}`)

	vb, em := f.Marshal(v)
	assert.NoError(em)

	logger.Info(string(vb))

	bv, eu := f.Unmarshal(b1)
	assert.NoError(eu)

	if string(b1) == string(vb) {
		assert.Equal(b1, vb)
	} else {
		assert.Equal(b2, vb)
	}
	proto.Equal(v, bv)
}

func TestProtojsonFormatUnMmarshalError(t *testing.T) {
	assert := assert.New(t)

	f := format.Protojson[*Message]()
	b := []byte(`{"topic":"test", brokenjson}`)

	bv, eu := f.Unmarshal(b)
	assert.Error(eu)
	assert.Equal(bv, &Message{})
}

func TestProtojsonFormatUnMmarshalEmpty(t *testing.T) {
	assert := assert.New(t)

	f := format.Protojson[*Message]()

	bv, eu := f.Unmarshal(nil)
	assert.NoError(eu)
	assert.Equal(bv, &Message{})
}
