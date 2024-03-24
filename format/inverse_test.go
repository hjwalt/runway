package format_test

import (
	"context"
	"testing"

	"github.com/hjwalt/runway/format"
	"github.com/hjwalt/runway/inverse"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestStoreFormat(t *testing.T) {
	assert := assert.New(t)

	c := inverse.NewContainer()
	f := format.Protobuf[*timestamppb.Timestamp]()

	format.RegisterFormat(c, f)
	val, err := format.RetrieveFormat[*timestamppb.Timestamp](c, context.Background())

	assert.NoError(err)
	assert.Equal(f, val)
}

func TestStoreMask(t *testing.T) {
	assert := assert.New(t)

	c := inverse.NewContainer()
	m := format.Base64Mask()

	format.RegisterMask(c, "base64", m)
	val, err := format.RetrieveMask(c, "base64", context.Background())

	assert.NoError(err)
	assert.Equal(m, val)
}
