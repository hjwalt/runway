package inverse_test

import (
	"context"
	"testing"

	"github.com/hjwalt/runway/inverse"
	"github.com/stretchr/testify/assert"
)

func TestGlobalResolveLast(t *testing.T) {
	cases := []struct {
		name          string
		reset         bool
		resolverAdder func()
		test          func(*testing.T, *assert.Assertions)
	}{
		{
			name:  "get last",
			reset: true,
			resolverAdder: func() {
				inverse.Register[string]("test-1", func(ctx context.Context) (string, error) { return "test-1", nil })
			},
			test: func(t *testing.T, a *assert.Assertions) {
				val, err := inverse.GetLast[string](context.Background(), "test-1")
				a.NoError(err)
				a.Equal("test-1", val)
			},
		},
		{
			name:  "get val",
			reset: true,
			resolverAdder: func() {
				inverse.RegisterInstance[string]("test-1", "test-1-instance")
			},
			test: func(t *testing.T, a *assert.Assertions) {
				val, err := inverse.GetLast[string](context.Background(), "test-1")
				a.NoError(err)
				a.Equal("test-1-instance", val)
			},
		},
		{
			name:  "get val cached",
			reset: false,
			resolverAdder: func() {
				inverse.RegisterInstance[string]("test-1", "test-1-instance-new")
			},
			test: func(t *testing.T, a *assert.Assertions) {
				val, err := inverse.GetLast[string](context.Background(), "test-1")
				a.NoError(err)
				a.Equal("test-1-instance", val)
			},
		},
		{
			name:  "get val released",
			reset: false,
			resolverAdder: func() {
				inverse.Release()
			},
			test: func(t *testing.T, a *assert.Assertions) {
				val, err := inverse.GetLast[string](context.Background(), "test-1")
				a.NoError(err)
				a.Equal("test-1-instance-new", val)
			},
		},
		{
			name:  "get vals",
			reset: true,
			resolverAdder: func() {
				inverse.RegisterInstances[string]("test-1", []string{"test-1-instance", "test-1-next"})
			},
			test: func(t *testing.T, a *assert.Assertions) {
				val, err := inverse.GetLast[string](context.Background(), "test-1")
				a.NoError(err)
				a.Equal("test-1-next", val)
			},
		},
		{
			name:  "get function",
			reset: true,
			resolverAdder: func() {
				inverse.Register("test-1", func(ctx context.Context) (any, error) {
					return inverse.Injector[string](func(ctx context.Context) (string, error) { return "test-1", nil }), nil
				})
			},
			test: func(t *testing.T, a *assert.Assertions) {
				val, err := inverse.GetLast[inverse.Injector[string]](context.Background(), "test-1")
				a.NoError(err)
				a.NotNil(val)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.reset {
				inverse.Reset()
			}
			tc.resolverAdder()
			tc.test(t, assert.New(t))
		})
	}
}

func TestGlobalResolveAll(t *testing.T) {
	cases := []struct {
		name          string
		reset         bool
		resolverAdder func()
		test          func(*testing.T, *assert.Assertions)
	}{
		{
			name:  "get all",
			reset: true,
			resolverAdder: func() {
				inverse.Register("test-1", func(ctx context.Context) (any, error) { return "test-1-a", nil })
				inverse.Register("test-1", func(ctx context.Context) (any, error) { return "test-1-b", nil })
			},
			test: func(t *testing.T, a *assert.Assertions) {
				val, err := inverse.GetAll[string](context.Background(), "test-1")
				a.NoError(err)
				a.Equal(2, len(val))
				a.Equal("test-1-a", val[0])
				a.Equal("test-1-b", val[1])
			},
		},
		{
			name:  "get vals",
			reset: true,
			resolverAdder: func() {
				inverse.RegisterInstances[string]("test-1", []string{"test-1-instance", "test-1-next"})
			},
			test: func(t *testing.T, a *assert.Assertions) {
				val, err := inverse.GetAll[string](context.Background(), "test-1")
				a.NoError(err)
				a.Equal(2, len(val))
				a.Equal("test-1-instance", val[0])
				a.Equal("test-1-next", val[1])
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.reset {
				inverse.Reset()
			}
			tc.resolverAdder()
			tc.test(t, assert.New(t))
		})
	}
}
