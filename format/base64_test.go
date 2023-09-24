package format_test

import (
	"testing"

	"github.com/hjwalt/runway/format"
	"github.com/stretchr/testify/assert"
)

func TestBase64Mask(t *testing.T) {
	assert := assert.New(t)

	base64MaskFormat := format.Base64String()

	assert.Equal("", base64MaskFormat.Default())
	assert.Equal([]byte{}, format.Base64().Default())

	masked, err := base64MaskFormat.Marshal("test")
	assert.NoError(err)
	assert.Equal("dGVzdA==", string(masked))

	unmasked, err := base64MaskFormat.Unmarshal(masked)

	assert.NoError(err)
	assert.Equal("test", unmasked)
}
