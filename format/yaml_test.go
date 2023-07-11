package format_test

import (
	"testing"

	"github.com/hjwalt/runway/format"
	"github.com/stretchr/testify/assert"
)

type TestYamlStruct struct {
	Name  string `yaml:"name"`
	Value int64  `yaml:"value"`
}

func TestYamlFormat(t *testing.T) {
	assert := assert.New(t)

	f := format.Yaml[TestYamlStruct]()
	v := TestYamlStruct{
		Name:  "test",
		Value: 1234567890,
	}
	b := []byte(`name: test
value: 1234567890
`)

	vb, em := f.Marshal(v)
	assert.NoError(em)

	bv, eu := f.Unmarshal(b)
	assert.NoError(eu)

	assert.Equal(b, vb)
	assert.Equal(v, bv)
}

func TestYamlFormatPointer(t *testing.T) {
	assert := assert.New(t)

	f := format.Yaml[*TestYamlStruct]()
	v := &TestYamlStruct{
		Name:  "test",
		Value: 1234567890,
	}
	b := []byte(`name: test
value: 1234567890
`)

	vb, em := f.Marshal(v)
	assert.NoError(em)

	bv, eu := f.Unmarshal(b)
	assert.NoError(eu)

	assert.Equal(b, vb)
	assert.Equal(v, bv)
}

func TestYamlFormatEmptyValue(t *testing.T) {
	assert := assert.New(t)

	f := format.Yaml[*TestYamlStruct]()

	vb, em := f.Marshal(nil)
	assert.NoError(em)
	assert.Equal([]byte{}, vb)

	bv, eu := f.Unmarshal(nil)
	assert.NoError(eu)
	assert.Equal(&TestYamlStruct{}, bv)
}
