package format_test

import (
	"testing"

	"github.com/hjwalt/runway/format"
	"github.com/stretchr/testify/assert"
)

func TestBase64Default(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("", format.Base64().Default())
}
