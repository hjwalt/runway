package runtime_test

import (
	"context"
	"testing"

	"github.com/hjwalt/runway/inverse"
	"github.com/hjwalt/runway/runtime"
	"github.com/stretchr/testify/assert"
)

func TestGetLastResolveWithConfigurationOptional(t *testing.T) {
	assert := assert.New(t)

	container := inverse.NewContainer()

	resolver := runtime.NewResolver[*runtime.HttpRunnable, runtime.Runtime](
		"test-1",
		container,
		false,
		runtime.NewHttp,
	)

	resolver.Register()

	httpRunnable, err := inverse.GenericGetLast[runtime.Runtime](container, context.Background(), "test-1")

	assert.NoError(err)
	assert.NotNil(httpRunnable)
}

func TestGetLastResolveWithConfigurationOptionalWithConfigError(t *testing.T) {
	assert := assert.New(t)

	container := inverse.NewContainer()

	resolver := runtime.NewResolver[*runtime.HttpRunnable, runtime.Runtime](
		"test-1",
		container,
		false,
		runtime.NewHttp,
	)

	container.AddVal(runtime.QualifierConfig("test-1"), "test")
	resolver.Register()

	httpRunnable, err := inverse.GenericGetLast[runtime.Runtime](container, context.Background(), "test-1")

	assert.Error(err)
	assert.Nil(httpRunnable)
}

func TestGetLastResolveWithConfigurationRequiredWithConfig(t *testing.T) {
	assert := assert.New(t)

	container := inverse.NewContainer()

	resolver := runtime.NewResolver[*runtime.HttpRunnable, runtime.Runtime](
		"test-1",
		container,
		true,
		runtime.NewHttp,
	)

	resolver.AddConfig(func(ctx context.Context) (runtime.Configuration[*runtime.HttpRunnable], error) {
		return runtime.HttpWithPort(8080), nil
	})
	resolver.Register()

	httpRunnable, err := inverse.GenericGetLast[runtime.Runtime](container, context.Background(), "test-1")

	assert.NoError(err)
	assert.NotNil(httpRunnable)
}

func TestGetLastResolveWithConfigurationRequiredWithConfigVal(t *testing.T) {
	assert := assert.New(t)

	container := inverse.NewContainer()

	resolver := runtime.NewResolver[*runtime.HttpRunnable, runtime.Runtime](
		"test-1",
		container,
		true,
		runtime.NewHttp,
	)

	resolver.AddConfigVal(runtime.HttpWithPort(8080))
	resolver.Register()

	httpRunnable, err := inverse.GenericGetLast[runtime.Runtime](container, context.Background(), "test-1")

	assert.NoError(err)
	assert.NotNil(httpRunnable)
}

func TestGetLastResolveWithConfigurationRequiredShouldFail(t *testing.T) {
	assert := assert.New(t)
	container := inverse.NewContainer()

	resolver := runtime.NewResolver[*runtime.HttpRunnable, runtime.Runtime](
		"test-1",
		container,
		true,
		runtime.NewHttp,
	)

	resolver.Register()

	httpRunnable, err := inverse.GenericGetLast[runtime.Runtime](container, context.Background(), "test-1")

	assert.Error(err)
	assert.ErrorIs(err, inverse.ErrInverseResolverMissing)
	assert.Nil(httpRunnable)
}
