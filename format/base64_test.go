package format_test

import (
	"testing"

	"github.com/hjwalt/runway/format"
	"github.com/stretchr/testify/assert"
)

func TestBase64(t *testing.T) {
	assert := assert.New(t)

	base64Format := format.Base64()

	assert.Equal("", base64Format.Default())

	masked, err := base64Format.Marshal("dGVzdA==")
	assert.NoError(err)
	assert.Equal("test", string(masked))

	unmasked, err := base64Format.Unmarshal(masked)

	assert.NoError(err)
	assert.Equal("dGVzdA==", unmasked)
}

func TestBase64Mask(t *testing.T) {
	assert := assert.New(t)

	base64MaskFormat := format.Masked(
		format.Base64Mask(),
		format.String(),
	)

	assert.Equal("", base64MaskFormat.Default())

	masked, err := base64MaskFormat.Marshal("test")
	assert.NoError(err)
	assert.Equal("dGVzdA==", string(masked))

	unmasked, err := base64MaskFormat.Unmarshal(masked)

	assert.NoError(err)
	assert.Equal("test", unmasked)
}
