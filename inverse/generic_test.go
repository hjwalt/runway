package inverse_test

import (
	"context"
	"testing"

	"github.com/hjwalt/runway/inverse"
	"github.com/stretchr/testify/assert"
)

func TestContainerGenericResolveLast(t *testing.T) {
	cases := []struct {
		name          string
		reset         bool
		resolverAdder func(inverse.Container)
		test          func(inverse.Container, *testing.T, *assert.Assertions)
	}{
		{
			name:  "get last",
			reset: false,
			resolverAdder: func(c inverse.Container) {
				c.Add("test-1", func(ctx context.Context, ci inverse.Container) (any, error) { return "test-1", nil })
			},
			test: func(c inverse.Container, t *testing.T, a *assert.Assertions) {
				val, err := inverse.GenericGetLast[string](c, context.Background(), "test-1")
				a.NoError(err)
				a.Equal("test-1", val)
			},
		},
		{
			name:  "get last cast error critical",
			reset: false,
			resolverAdder: func(c inverse.Container) {
			},
			test: func(c inverse.Container, t *testing.T, a *assert.Assertions) {
				val, err := inverse.GenericGetLast[int64](c, context.Background(), "test-1")
				a.ErrorIs(err, inverse.ErrInverseCastError)
				a.Equal(int64(0), val)
				valstr, err := inverse.GenericGetLast[string](c, context.Background(), "test-1")
				a.ErrorIs(err, inverse.ErrInverseCastError)
				a.Equal("", valstr)
			},
		},
		{
			name:  "critical error reset",
			reset: true,
			resolverAdder: func(c inverse.Container) {
				c.Add("test-1", func(ctx context.Context, ci inverse.Container) (any, error) { return "test-1", nil })
			},
			test: func(c inverse.Container, t *testing.T, a *assert.Assertions) {
				val, err := inverse.GenericGetLast[string](c, context.Background(), "test-1")
				a.NoError(err)
				a.Equal("test-1", val)
			},
		},
		{
			name:  "get last val",
			reset: true,
			resolverAdder: func(c inverse.Container) {
				inverse.GenericAddVal(c, "test-1", "test-2")
			},
			test: func(c inverse.Container, t *testing.T, a *assert.Assertions) {
				val, err := inverse.GenericGetLast[string](c, context.Background(), "test-1")
				a.NoError(err)
				a.Equal("test-2", val)
			},
		},
	}

	c := inverse.NewContainer()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.reset {
				c.Reset()
			}
			tc.resolverAdder(c)
			tc.test(c, t, assert.New(t))
		})
	}
}

func TestContainerGenericResolveAll(t *testing.T) {
	cases := []struct {
		name          string
		reset         bool
		resolverAdder func(inverse.Container)
		test          func(inverse.Container, *testing.T, *assert.Assertions)
	}{
		{
			name:  "get all",
			reset: false,
			resolverAdder: func(c inverse.Container) {
				c.Add("test-1", func(ctx context.Context, ci inverse.Container) (any, error) { return "test-1-a", nil })
				c.Add("test-1", func(ctx context.Context, ci inverse.Container) (any, error) { return "test-1-b", nil })
			},
			test: func(c inverse.Container, t *testing.T, a *assert.Assertions) {
				val, err := inverse.GenericGetAll[string](c, context.Background(), "test-1")
				a.NoError(err)
				a.Equal(2, len(val))
				a.Equal("test-1-a", val[0])
				a.Equal("test-1-b", val[1])
			},
		},
		{
			name:  "get all cast error critical",
			reset: false,
			resolverAdder: func(c inverse.Container) {
			},
			test: func(c inverse.Container, t *testing.T, a *assert.Assertions) {
				val, err := inverse.GenericGetAll[int64](c, context.Background(), "test-1")
				a.ErrorIs(err, inverse.ErrInverseCastError)
				a.Equal(0, len(val))
				valstr, err := inverse.GenericGetAll[string](c, context.Background(), "test-1")
				a.ErrorIs(err, inverse.ErrInverseCastError)
				a.Equal(0, len(valstr))
			},
		},
		{
			name:  "critical error reset",
			reset: true,
			resolverAdder: func(c inverse.Container) {
				c.Add("test-1", func(ctx context.Context, ci inverse.Container) (any, error) { return "test-1-a", nil })
				c.Add("test-1", func(ctx context.Context, ci inverse.Container) (any, error) { return "test-1-b", nil })
			},
			test: func(c inverse.Container, t *testing.T, a *assert.Assertions) {
				val, err := inverse.GenericGetAll[string](c, context.Background(), "test-1")
				a.NoError(err)
				a.Equal(2, len(val))
				a.Equal("test-1-a", val[0])
				a.Equal("test-1-b", val[1])
			},
		},
	}

	c := inverse.NewContainer()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.reset {
				c.Reset()
			}
			tc.resolverAdder(c)
			tc.test(c, t, assert.New(t))
		})
	}
}
