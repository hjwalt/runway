package format_test

import (
	"testing"

	"github.com/hjwalt/runway/format"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestProtobufFormat(t *testing.T) {
	assert := assert.New(t)

	base64format := format.Base64()

	f := format.Protobuf[*timestamppb.Timestamp]()
	v := &timestamppb.Timestamp{
		Seconds: 1234567890,
		Nanos:   123456789,
	}
	b64 := []byte("CNKF2MwEEJWa7zo=")

	b, e64 := base64format.Unmarshal(b64)
	assert.NoError(e64)

	vb, em := f.Marshal(v)
	assert.NoError(em)

	bv, eu := f.Unmarshal(b)
	assert.NoError(eu)

	vb64, e64 := base64format.Marshal(vb)
	assert.NoError(e64)

	assert.Equal(b64, vb64)
	proto.Equal(v, bv)
}

func TestProtobufFormatUnmarshalFailure(t *testing.T) {
	assert := assert.New(t)

	f := format.Protobuf[*timestamppb.Timestamp]()

	bv, eu := f.Unmarshal([]byte("test"))
	assert.Error(eu)
	assert.NotNil(bv)
}
