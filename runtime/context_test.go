package runtime_test

import (
	"context"
	"testing"

	"github.com/hjwalt/runway/runtime"
	"github.com/stretchr/testify/assert"
)

func TestContextKey(t *testing.T) {
	assert := assert.New(t)

	keyedContext := context.WithValue(context.Background(), runtime.Context("test"), "test-value")
	keyedContext = context.WithValue(keyedContext, runtime.Context("test-2"), "test-value-2")

	assert.Equal("test-value", keyedContext.Value(runtime.Context("test")))
	assert.Equal("test-value-2", keyedContext.Value(runtime.Context("test-2")))
}

func TestContextValueMap(t *testing.T) {
	assert := assert.New(t)

	valueMap := map[string]any{
		"test":   "test-value-map",
		"test-2": "test-value-map-2",
	}

	keyedContext := runtime.ContextWithValueMap(context.Background(), valueMap)

	assert.Equal("test-value-map", keyedContext.Value(runtime.Context("test")))
	assert.Equal("test-value-map-2", keyedContext.Value(runtime.Context("test-2")))
}

func TestContextValue(t *testing.T) {
	assert := assert.New(t)

	valueMap := map[string]any{
		"string":   "value-string",
		"int32str": "32",
		"int32":    int32(32),
		"int64str": "64",
		"int64":    int64(64),
		"boolStr":  "true",
		"bool":     true,
	}

	keyedContext := runtime.ContextWithValueMap(context.Background(), valueMap)

	assert.Equal("value-string", runtime.ContextValueString(keyedContext, "string"))
	assert.Equal(int32(32), runtime.ContextValueInt32(keyedContext, "int32str"))
	assert.Equal(int32(32), runtime.ContextValueInt32(keyedContext, "int32"))
	assert.Equal(int64(64), runtime.ContextValueInt64(keyedContext, "int64str"))
	assert.Equal(int64(64), runtime.ContextValueInt64(keyedContext, "int64"))
	assert.Equal(true, runtime.ContextValueBool(keyedContext, "boolStr"))
	assert.Equal(true, runtime.ContextValueBool(keyedContext, "bool"))

	assert.Equal("", runtime.ContextValueString(keyedContext, "missing"))
	assert.Equal(int32(0), runtime.ContextValueInt32(keyedContext, "missing"))
	assert.Equal(int64(0), runtime.ContextValueInt64(keyedContext, "missing"))
	assert.Equal(false, runtime.ContextValueBool(keyedContext, "missing"))
}
