package inverse_test

import (
	"context"
	"testing"

	"github.com/hjwalt/runway/inverse"
	"github.com/stretchr/testify/assert"
)

func TestContainerResolveLast(t *testing.T) {
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
				c.Add("test-1", func(ctx context.Context) (any, error) { return "test-1", nil })
			},
			test: func(c inverse.Container, t *testing.T, a *assert.Assertions) {
				val, err := c.Get(context.Background(), "test-1")
				a.NoError(err)
				a.Equal("test-1", val)
			},
		},
		{
			name:  "get last from cache",
			reset: false,
			resolverAdder: func(c inverse.Container) {
				c.Add("test-1", func(ctx context.Context) (any, error) { return "test-1-new", nil })
			},
			test: func(c inverse.Container, t *testing.T, a *assert.Assertions) {
				val, err := c.Get(context.Background(), "test-1")
				a.NoError(err)
				a.Equal("test-1", val)
			},
		},
		{
			name:  "get new last after clear",
			reset: false,
			resolverAdder: func(c inverse.Container) {
				c.Clear()
			},
			test: func(c inverse.Container, t *testing.T, a *assert.Assertions) {
				val, err := c.Get(context.Background(), "test-1")
				a.NoError(err)
				a.Equal("test-1-new", val)
			},
		},
		{
			name:  "get val",
			reset: true,
			resolverAdder: func(c inverse.Container) {
				c.AddVal("test-1", "test-1-val")
			},
			test: func(c inverse.Container, t *testing.T, a *assert.Assertions) {
				val, err := c.Get(context.Background(), "test-1")
				a.NoError(err)
				a.Equal("test-1-val", val)
			},
		},
		{
			name:  "get last multiple added",
			reset: true,
			resolverAdder: func(c inverse.Container) {
				c.Add("test-1", func(ctx context.Context) (any, error) { return "test-1-a", nil })
				c.Add("test-1", func(ctx context.Context) (any, error) { return "test-1-b", nil })
			},
			test: func(c inverse.Container, t *testing.T, a *assert.Assertions) {
				val, err := c.Get(context.Background(), "test-1")
				a.NoError(err)
				a.Equal("test-1-b", val)
			},
		},
		{
			name:  "fail missing context",
			reset: true,
			resolverAdder: func(c inverse.Container) {
				c.Add("test-1", func(ctx context.Context) (any, error) { return "test-1-fail", nil })
			},
			test: func(c inverse.Container, t *testing.T, a *assert.Assertions) {
				val, err := c.Get(nil, "test-1")
				a.ErrorIs(err, inverse.ErrInverseNilContext)
				a.Nil(val)
				val, err = c.Get(context.Background(), "test-1")
				a.NoError(err)
				a.Equal("test-1-fail", val)
			},
		},
		{
			name:  "fail missing resolver",
			reset: true,
			resolverAdder: func(c inverse.Container) {
				c.Add("test-1", func(ctx context.Context) (any, error) { return "test-1-fail", nil })
			},
			test: func(c inverse.Container, t *testing.T, a *assert.Assertions) {
				val, err := c.Get(context.Background(), "test-2")
				a.ErrorIs(err, inverse.ErrInverseResolverMissing)
				a.Nil(val)
				val, err = c.Get(context.Background(), "test-1")
				a.NoError(err)
				a.Equal("test-1-fail", val)
			},
		},
		{
			name:  "fail resolve loop critical error",
			reset: true,
			resolverAdder: func(c inverse.Container) {
				c.Add("test-1", func(ctx context.Context) (any, error) {
					_, err := c.Get(ctx, "test-2")
					if err != nil {
						return "", err
					}
					return "test-1", nil
				})
				c.Add("test-2", func(ctx context.Context) (any, error) {
					_, err := c.Get(ctx, "test-1")
					if err != nil {
						return "", err
					}
					return "test-2", nil
				})
				c.Add("test-3", func(ctx context.Context) (any, error) { return "test-3-no-loop", nil })
			},
			test: func(c inverse.Container, t *testing.T, a *assert.Assertions) {
				val, err := c.Get(context.Background(), "test-2")
				a.ErrorIs(err, inverse.ErrInverseResolveLoop)
				a.Nil(val)
				val, err = c.Get(context.Background(), "test-3")
				a.ErrorIs(err, inverse.ErrInverseResolveLoop)
				a.Nil(val)
			},
		},
		{
			name:  "critical error add does not do anything",
			reset: false,
			resolverAdder: func(c inverse.Container) {
				c.Add("test-critical", func(ctx context.Context) (any, error) { return "test-critical", nil })
			},
			test: func(c inverse.Container, t *testing.T, a *assert.Assertions) {
				val, err := c.Get(context.Background(), "test-critical")
				a.ErrorIs(err, inverse.ErrInverseResolveLoop)
				a.Nil(val)
			},
		},
		{
			name:  "critical error add val does not do anything",
			reset: false,
			resolverAdder: func(c inverse.Container) {
				c.AddVal("test-critical", "test-critical")
			},
			test: func(c inverse.Container, t *testing.T, a *assert.Assertions) {
				val, err := c.Get(context.Background(), "test-critical")
				a.ErrorIs(err, inverse.ErrInverseResolveLoop)
				a.Nil(val)
			},
		},
		{
			name:  "critical error cleaned on reset",
			reset: true,
			resolverAdder: func(c inverse.Container) {
				c.Add("test-1", func(ctx context.Context) (any, error) { return "test-1", nil })
			},
			test: func(c inverse.Container, t *testing.T, a *assert.Assertions) {
				val, err := c.Get(context.Background(), "test-1")
				a.NoError(err)
				a.Equal("test-1", val)
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

func TestContainerResolveAll(t *testing.T) {
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
				c.Add("test-1", func(ctx context.Context) (any, error) { return "test-1-a", nil })
				c.Add("test-1", func(ctx context.Context) (any, error) { return "test-1-b", nil })
			},
			test: func(c inverse.Container, t *testing.T, a *assert.Assertions) {
				val, err := c.GetAll(context.Background(), "test-1")
				a.NoError(err)
				a.Equal(2, len(val))
				a.Equal("test-1-a", val[0])
				a.Equal("test-1-b", val[1])
			},
		},
		{
			name:  "get last from cache",
			reset: false,
			resolverAdder: func(c inverse.Container) {
				c.Add("test-1", func(ctx context.Context) (any, error) { return "test-1-new", nil })
			},
			test: func(c inverse.Container, t *testing.T, a *assert.Assertions) {
				val, err := c.GetAll(context.Background(), "test-1")
				a.NoError(err)
				a.Equal(2, len(val))
				a.Equal("test-1-a", val[0])
				a.Equal("test-1-b", val[1])
			},
		},
		{
			name:  "fail missing context",
			reset: true,
			resolverAdder: func(c inverse.Container) {
				c.Add("test-1", func(ctx context.Context) (any, error) { return "test-1-a-fail", nil })
				c.Add("test-1", func(ctx context.Context) (any, error) { return "test-1-b-fail", nil })
			},
			test: func(c inverse.Container, t *testing.T, a *assert.Assertions) {
				val, err := c.GetAll(nil, "test-1")
				a.ErrorIs(err, inverse.ErrInverseNilContext)
				a.NotNil(val)
				a.Equal(0, len(val))
				val, err = c.GetAll(context.Background(), "test-1")
				a.NoError(err)
				a.Equal(2, len(val))
				a.Equal("test-1-a-fail", val[0])
				a.Equal("test-1-b-fail", val[1])
			},
		},
		{
			name:  "fail missing resolver",
			reset: true,
			resolverAdder: func(c inverse.Container) {
				c.Add("test-1", func(ctx context.Context) (any, error) { return "test-1-a-fail", nil })
				c.Add("test-1", func(ctx context.Context) (any, error) { return "test-1-b-fail", nil })
			},
			test: func(c inverse.Container, t *testing.T, a *assert.Assertions) {
				val, err := c.GetAll(context.Background(), "test-2")
				a.ErrorIs(err, inverse.ErrInverseResolverMissing)
				a.NotNil(val)
				a.Equal(0, len(val))
				val, err = c.GetAll(context.Background(), "test-1")
				a.NoError(err)
				a.Equal(2, len(val))
				a.Equal("test-1-a-fail", val[0])
				a.Equal("test-1-b-fail", val[1])
			},
		},
		{
			name:  "fail resolve loop critical error",
			reset: true,
			resolverAdder: func(c inverse.Container) {
				c.Add("test-1", func(ctx context.Context) (any, error) {
					_, err := c.Get(ctx, "test-2")
					if err != nil {
						return "", err
					}
					return "test-1", nil
				})
				c.Add("test-2", func(ctx context.Context) (any, error) {
					_, err := c.Get(ctx, "test-1")
					if err != nil {
						return "", err
					}
					return "test-2", nil
				})
				c.Add("test-3", func(ctx context.Context) (any, error) { return "test-3-no-loop", nil })
			},
			test: func(c inverse.Container, t *testing.T, a *assert.Assertions) {
				val, err := c.GetAll(context.Background(), "test-2")
				a.ErrorIs(err, inverse.ErrInverseResolveLoop)
				a.NotNil(val)
				a.Equal(0, len(val))
				val, err = c.GetAll(context.Background(), "test-3")
				a.ErrorIs(err, inverse.ErrInverseResolveLoop)
				a.NotNil(val)
				a.Equal(0, len(val))
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
