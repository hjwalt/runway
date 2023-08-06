package inverse_test

import (
	"context"
	"testing"

	"github.com/hjwalt/runway/inverse"
	"github.com/stretchr/testify/assert"
)

func TestGetAllResolveOneLevel(t *testing.T) {
	assert := assert.New(t)
	inverse.Reset()

	inverse.Register("test-1", func(ctx context.Context) (any, error) { return "test-1", nil })
	inverse.Register("test-1", func(ctx context.Context) (string, error) { return "test-1-next", nil })

	val, err := inverse.GetAll[string](context.Background(), "test-1")

	assert.NoError(err)
	assert.Equal(2, len(val))
	assert.Equal("test-1", val[0])
	assert.Equal("test-1-next", val[1])
}

func TestGetAllResolveTwoLevel(t *testing.T) {
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
	inverse.Register("test-2", func(ctx context.Context) (any, error) { return "test-2", nil })

	val, err := inverse.GetAll[string](context.Background(), "test-2")

	assert.NoError(err)
	assert.Equal(2, len(val))
	assert.Equal("test-2-test-1", val[0])
	assert.Equal("test-2", val[1])
}

func TestGetAllResolvePartialForCastingError(t *testing.T) {
	assert := assert.New(t)
	inverse.Reset()

	inverse.Register("test-1", func(ctx context.Context) (any, error) { return int64(1), nil })
	inverse.Register("test-1", func(ctx context.Context) (any, error) { return "test-1", nil })
	inverse.Register("test-1", func(ctx context.Context) (any, error) { return int64(2), nil })

	val, err := inverse.GetAll[int64](context.Background(), "test-1")

	assert.Error(err)
	assert.Equal(inverse.ErrorCastingFailure("test-1"), err)
	assert.Equal(2, len(val))
	assert.Equal(int64(1), val[0])
	assert.Equal(int64(2), val[1])
}

func TestGetAllFailForMissingContext(t *testing.T) {
	assert := assert.New(t)
	inverse.Reset()

	inverse.Register("test-1", func(ctx context.Context) (any, error) { return "test-1", nil })

	val, err := inverse.GetAll[string](nil, "test-1")

	assert.Error(err)
	assert.Equal(inverse.ErrorNilContext("test-1"), err)
	assert.Equal(0, len(val))
}

func TestGetAllFailForMissingQualifier(t *testing.T) {
	assert := assert.New(t)
	inverse.Reset()

	inverse.Register("test-1", func(ctx context.Context) (any, error) { return "test-1", nil })

	val, err := inverse.GetAll[string](context.Background(), "test-2")

	assert.Error(err)
	assert.Equal(inverse.ErrorNotInjected("test-2"), err)
	assert.Equal(0, len(val))
}

func TestGetAllFailForResolveLoop(t *testing.T) {
	assert := assert.New(t)
	inverse.Reset()

	inverse.Register("test-1", func(ctx context.Context) (any, error) {
		_, err := inverse.GetAll[string](ctx, "test-2")
		if err != nil {
			return "", err
		}
		return "test-1", nil
	})
	inverse.Register("test-2", func(ctx context.Context) (any, error) {
		_, err := inverse.GetAll[string](ctx, "test-1")
		if err != nil {
			return "", err
		}
		return "test-2", nil
	})

	val, err := inverse.GetAll[string](context.Background(), "test-1")

	assert.Error(err)
	assert.Equal(inverse.ErrorResolveLoop("test-1"), err)
	assert.Equal(0, len(val))
}
