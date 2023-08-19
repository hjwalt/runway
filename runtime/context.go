package runtime

import (
	"context"

	"github.com/hjwalt/runway/reflect"
)

type ContextKey struct {
	Name string
}

func Context(name string) ContextKey {
	return ContextKey{Name: name}
}

func ContextWithValueMap(ctx context.Context, valueMap map[string]any) context.Context {
	updatedContext := ctx
	for k, v := range valueMap {
		updatedContext = context.WithValue(updatedContext, Context(k), v)
	}
	return updatedContext
}

func ContextValueString(ctx context.Context, name string) string {
	return reflect.GetString(ctx.Value(Context(name)))
}

func ContextValueInt32(ctx context.Context, name string) int32 {
	return reflect.GetInt32(ctx.Value(Context(name)))
}

func ContextValueInt64(ctx context.Context, name string) int64 {
	return reflect.GetInt64(ctx.Value(Context(name)))
}

func ContextValueBool(ctx context.Context, name string) bool {
	return reflect.GetBool(ctx.Value(Context(name)))
}
