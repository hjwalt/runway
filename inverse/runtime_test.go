package inverse_test

import (
	"context"
	"testing"

	"github.com/hjwalt/runway/inverse"
	"github.com/hjwalt/runway/runtime"
	"github.com/stretchr/testify/assert"
)

func TestGetLastResolveWithConfigurationOptional(t *testing.T) {
	assert := assert.New(t)
	inverse.Reset()

	inverse.RegisterWithConfigurationOptional("test-1", "test-1-configuration", runtime.NewHttp)

	httpRunnable, err := inverse.GetLast[runtime.Runtime](context.Background(), "test-1")

	assert.NoError(err)
	assert.NotNil(httpRunnable)
}

func TestGetLastResolveWithConfigurationOptionalWithConfigError(t *testing.T) {
	assert := assert.New(t)
	inverse.Reset()

	inverse.RegisterInstance[string]("test-1-configuration", "test")
	inverse.RegisterWithConfigurationOptional("test-1", "test-1-configuration", runtime.NewHttp)

	httpRunnable, err := inverse.GetLast[runtime.Runtime](context.Background(), "test-1")

	assert.Error(err)
	assert.Nil(httpRunnable)
}

func TestGetLastResolveWithConfigurationOptionalWithConfig(t *testing.T) {
	assert := assert.New(t)
	inverse.Reset()

	inverse.RegisterConfiguration[*runtime.HttpRunnable]("test-1-configuration", runtime.HttpWithPort(8080))
	inverse.RegisterWithConfigurationOptional("test-1", "test-1-configuration", runtime.NewHttp)

	httpRunnable, err := inverse.GetLast[runtime.Runtime](context.Background(), "test-1")

	assert.NoError(err)
	assert.NotNil(httpRunnable)
}

func TestGetLastResolveWithConfigurationRequired(t *testing.T) {
	assert := assert.New(t)
	inverse.Reset()

	inverse.RegisterConfigurations[*runtime.HttpRunnable]("test-1-configuration", []runtime.Configuration[*runtime.HttpRunnable]{runtime.HttpWithPort(8080)})
	inverse.RegisterWithConfigurationRequired("test-1", "test-1-configuration", runtime.NewHttp)

	httpRunnable, err := inverse.GetLast[runtime.Runtime](context.Background(), "test-1")

	assert.NoError(err)
	assert.NotNil(httpRunnable)
}

func TestGetLastResolveWithConfigurationRequiredShouldFail(t *testing.T) {
	assert := assert.New(t)
	inverse.Reset()

	inverse.RegisterWithConfigurationRequired("test-1", "test-1-configuration", runtime.NewHttp)

	httpRunnable, err := inverse.GetLast[runtime.Runtime](context.Background(), "test-1")

	assert.Error(err)
	assert.ErrorIs(err, inverse.ErrInverseResolverMissing)
	assert.Nil(httpRunnable)
}
