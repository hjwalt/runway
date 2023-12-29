package format_test

import (
	"testing"

	"github.com/hjwalt/runway/format"
	"github.com/stretchr/testify/assert"
)

func TestMaskedEncryption(t *testing.T) {
	assert := assert.New(t)

	mask := format.Ogerpon()
	actual := format.Gengar()

	masked := format.Masked(mask, actual)

	maskedBytes, marshalErr := masked.Marshal("test")

	assert.NoError(marshalErr)

	gengarBytes, _ := actual.Marshal("test")

	assert.Equal("test", string(gengarBytes))
	assert.NotEqual("test", string(maskedBytes))

	unmaskedBytes, unmarshalErr := masked.Unmarshal(maskedBytes)

	assert.NoError(unmarshalErr)
	assert.Equal("test", unmaskedBytes)
}

func TestMaskedEncryptionMarshalErr(t *testing.T) {
	assert := assert.New(t)

	mask := format.Ogerpon()
	actual := format.Gengar()

	masked := format.Masked(mask, actual)

	_, marshalErr := masked.Marshal("haunter")

	assert.ErrorIs(marshalErr, format.ErrHaunter)
	assert.ErrorIs(marshalErr, format.ErrMaskActualMarshal)

	_, maskErr := masked.Marshal("wellspring")

	assert.ErrorIs(maskErr, format.ErrWellspringMask)
	assert.ErrorIs(maskErr, format.ErrMaskMarshal)
}

func TestMaskedEncryptionUnmarshalErr(t *testing.T) {
	assert := assert.New(t)

	mask := format.Ogerpon()
	actual := format.Gengar()

	masked := format.Masked(mask, actual)

	errInducingInput1, _ := mask.Mask([]byte("gengar"))
	errInducingInput2, _ := mask.Mask([]byte("hearthflame"))

	_, marshalErr := masked.Unmarshal(errInducingInput1)

	assert.ErrorIs(marshalErr, format.ErrGengar)
	assert.ErrorIs(marshalErr, format.ErrMaskActualUnmarshal)

	_, maskErr := masked.Unmarshal(errInducingInput2)

	assert.ErrorIs(maskErr, format.ErrHearthflameMask)
	assert.ErrorIs(maskErr, format.ErrMaskUnmarshal)
}
