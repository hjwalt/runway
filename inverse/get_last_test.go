package inverse_test

import (
	"context"
	"testing"

	"github.com/hjwalt/runway/format"
	"github.com/hjwalt/runway/inverse"
	"github.com/hjwalt/runway/runtime"
	"github.com/stretchr/testify/assert"
)

func TestGetLastResolveOneLevel(t *testing.T) {
	assert := assert.New(t)
	inverse.Reset()

	inverse.Register("test-1", func(ctx context.Context) (string, error) { return "test-1", nil })

	val, err := inverse.GetLast[string](context.Background(), "test-1")

	assert.NoError(err)
	assert.Equal("test-1", val)
}

func TestGetLastResolveTypeCasting(t *testing.T) {
	assert := assert.New(t)
	inverse.Reset()

	inverse.Register("strfmt", func(ctx context.Context) (any, error) { return format.String(), nil })

	val, err := inverse.GetLast[format.Format[string]](context.Background(), "strfmt")

	assert.NoError(err)
	assert.Equal(format.String(), val)
}

func TestGetLastResolveLastForQualifier(t *testing.T) {
	assert := assert.New(t)
	inverse.Reset()

	inverse.Register("test-1", func(ctx context.Context) (any, error) { return "test-1", nil })
	inverse.Register("test-1", func(ctx context.Context) (any, error) { return "test-1-last", nil })

	val, err := inverse.GetLast[string](context.Background(), "test-1")

	assert.NoError(err)
	assert.Equal("test-1-last", val)
}

func TestGetLastResolveTwoLevelOfInjection(t *testing.T) {
	assert := assert.New(t)
	inverse.Reset()

	inverse.Register("test-1", func(ctx context.Context) (any, error) { return "test-1", nil })
	inverse.Register("test-2", func(ctx context.Context) (any, error) {
		val, err := inverse.GetLast[string](ctx, "test-1")
		if err != nil {
			return "", err
		}
		return "test-2-" + val, nil
	})

	val, err := inverse.GetLast[string](context.Background(), "test-2")

	assert.NoError(err)
	assert.Equal("test-2-test-1", val)
}

func TestGetLastResolveCacheIfExist(t *testing.T) {
	assert := assert.New(t)
	inverse.Reset()

	inverse.Register("test-1", func(ctx context.Context) (any, error) { return "test-1", nil })
	inverse.Register("test-1", func(ctx context.Context) (any, error) { return "test-1-last", nil })
	val, err := inverse.GetLast[string](context.Background(), "test-1")

	assert.NoError(err)
	assert.Equal("test-1-last", val)

	inverse.Register("test-1", func(ctx context.Context) (any, error) { return "test-1-updated", nil })

	val, err = inverse.GetLast[string](context.Background(), "test-1")

	assert.NoError(err)
	assert.Equal("test-1-last", val)
}

func TestGetLastResolveFunction(t *testing.T) {
	assert := assert.New(t)
	inverse.Reset()

	inverse.Register("test-1", func(ctx context.Context) (any, error) {
		return inverse.Injector[string](func(ctx context.Context) (string, error) { return "test-1", nil }), nil
	})

	val, err := inverse.GetLast[inverse.Injector[string]](context.Background(), "test-1")

	assert.NoError(err)
	assert.NotNil(val)

	valRes, valErr := val(context.Background())
	assert.NoError(valErr)
	assert.Equal("test-1", valRes)
}

func TestGetLastFailForMissingContext(t *testing.T) {
	assert := assert.New(t)
	inverse.Reset()

	inverse.Register("test-1", func(ctx context.Context) (any, error) { return "test-1", nil })

	val, err := inverse.GetLast[string](nil, "test-1")

	assert.Error(err)
	assert.Equal(inverse.ErrorNilContext("test-1"), err)
	assert.Equal("", val)
}

func TestGetLastFailForMissingQualifier(t *testing.T) {
	assert := assert.New(t)
	inverse.Reset()

	inverse.Register("test-1", func(ctx context.Context) (any, error) { return "test-1", nil })

	val, err := inverse.GetLast[string](context.Background(), "test-2")

	assert.Error(err)
	assert.Equal(inverse.ErrorNotInjected("test-2"), err)
	assert.Equal("", val)
}

func TestGetLastFailForCastingError(t *testing.T) {
	assert := assert.New(t)
	inverse.Reset()

	inverse.Register("test-1", func(ctx context.Context) (any, error) { return "test-1", nil })

	val, err := inverse.GetLast[int64](context.Background(), "test-1")

	assert.Error(err)
	assert.Equal(inverse.ErrorCastingFailure("test-1"), err)
	assert.Equal(int64(0), val)
}

func TestGetLastFailForCacheCastingFailure(t *testing.T) {
	assert := assert.New(t)
	inverse.Reset()

	inverse.Register("test-1", func(ctx context.Context) (any, error) { return "test-1", nil })
	inverse.Register("test-1", func(ctx context.Context) (any, error) { return "test-1-last", nil })
	val, err := inverse.GetLast[string](context.Background(), "test-1")

	assert.NoError(err)
	assert.Equal("test-1-last", val)

	inverse.Register("test-1", func(ctx context.Context) (any, error) { return "test-1-updated", nil })

	valInt, err := inverse.GetLast[int64](context.Background(), "test-1")

	assert.Error(err)
	assert.Equal(inverse.ErrorCastingFailure("test-1"), err)
	assert.Equal(int64(0), valInt)
}

func TestGetLastFailForResolveLoop(t *testing.T) {
	assert := assert.New(t)
	inverse.Reset()

	inverse.Register("test-1", func(ctx context.Context) (any, error) {
		_, err := inverse.GetLast[string](ctx, "test-2")
		if err != nil {
			return "", err
		}
		return "test-1", nil
	})
	inverse.Register("test-2", func(ctx context.Context) (any, error) {
		_, err := inverse.GetLast[string](ctx, "test-1")
		if err != nil {
			return "", err
		}
		return "test-2", nil
	})

	val, err := inverse.GetLast[string](context.Background(), "test-1")

	assert.Error(err)
	assert.Equal(inverse.ErrorResolveLoop("test-1"), err)
	assert.Equal("", val)
}

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
	assert.ErrorIs(err, inverse.ErrNotInjected)
	assert.Nil(httpRunnable)
}
