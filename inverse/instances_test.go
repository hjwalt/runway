package inverse_test

import (
	"context"
	"testing"

	"github.com/hjwalt/runway/inverse"
	"github.com/stretchr/testify/assert"
)

func TestGetLastResolveInstance(t *testing.T) {
	assert := assert.New(t)
	inverse.Reset()

	inverse.RegisterInstance("test-1", "test-value")

	val, err := inverse.GetLast[string](context.Background(), "test-1")

	assert.NoError(err)
	assert.Equal("test-value", val)
}

func TestGetLastResolveInstancesLastValue(t *testing.T) {
	assert := assert.New(t)
	inverse.Reset()

	inverse.RegisterInstances("test-1", []string{"test-value", "test-value-last"})

	val, err := inverse.GetLast[string](context.Background(), "test-1")

	assert.NoError(err)
	assert.Equal("test-value-last", val)
}

func TestGetLastResolveEmptyInstancesLastValue(t *testing.T) {
	assert := assert.New(t)
	inverse.Reset()

	inverse.RegisterInstances("test-1", []string{})

	_, err := inverse.GetLast[string](context.Background(), "test-1")

	assert.Error(err)
	assert.ErrorIs(err, inverse.ErrNotInjected)
}
