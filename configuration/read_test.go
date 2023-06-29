package configuration_test

import (
	"testing"

	"github.com/hjwalt/runway/configuration"
	"github.com/hjwalt/runway/format"
	"github.com/stretchr/testify/assert"
)

type Config struct {
	Name   string `yaml:"name"`
	Number int64  `yaml:"number"`
}

func TestRead(t *testing.T) {
	assert := assert.New(t)

	c, err := configuration.Read("fixture/test.yaml", format.Yaml[Config]())

	assert.NoError(err)
	assert.Equal("pikachu", c.Name)
	assert.Equal(int64(1234), c.Number)
}

func TestReadFail(t *testing.T) {
	assert := assert.New(t)

	c, err := configuration.Read("fixture/nothere.yaml", format.Yaml[Config]())

	assert.ErrorIs(err, configuration.ErrReadFail)
	assert.Equal(Config{}, c)
}
