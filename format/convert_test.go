package format_test

import (
	"testing"

	"github.com/hjwalt/runway/format"
	"github.com/stretchr/testify/assert"
)

func TestConvert(t *testing.T) {
	assert := assert.New(t)

	fs := format.Json[TestJsonStruct]()
	ft := format.String()

	vs := TestJsonStruct{
		Name:  "test",
		Value: 1234567890,
	}
	et := `{"Name":"test","Value":1234567890}`

	vt, ferr := format.Convert(vs, fs, ft)
	assert.NoError(ferr)
	assert.Equal(et, vt)
}

func TestConvertUnmarshalError(t *testing.T) {
	assert := assert.New(t)

	fs := format.Gengar()
	ft := format.Gengar()

	vs := "ghastly"

	_, ferr := format.Convert(vs, fs, ft)
	assert.ErrorIs(ferr, format.ErrFormatConversionUnmarshal)
	assert.ErrorIs(ferr, format.ErrGhastly)
}

func TestConvertMarshalError(t *testing.T) {
	assert := assert.New(t)

	fs := format.Gengar()
	ft := format.Gengar()

	vs := "haunter"

	_, ferr := format.Convert(vs, fs, ft)
	assert.ErrorIs(ferr, format.ErrFormatConversionMarshal)
	assert.ErrorIs(ferr, format.ErrHaunter)
}
