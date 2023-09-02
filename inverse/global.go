package inverse

import (
	"context"
)

type Injector[T any] func(context.Context) (T, error)

var Global = NewContainer()

func Reset() {
	Global.Reset()
}

func Release() {
	Global.Clear()
}

func Register[T any](qualifier string, injector Injector[T]) {
	Global.Add(qualifier, func(ctx context.Context) (any, error) { return injector(ctx) })
}

func RegisterInstance[T any](qualifier string, instance T) {
	Global.AddVal(qualifier, instance)
}

func RegisterInstances[T any](qualifier string, instances []T) {
	for _, instance := range instances {
		Global.AddVal(qualifier, instance)
	}
}

func GetLast[T any](ctx context.Context, qualifier string) (T, error) {
	return GenericGetLast[T](Global, ctx, qualifier)
}

func GetAll[T any](ctx context.Context, qualifier string) ([]T, error) {
	return GenericGetAll[T](Global, ctx, qualifier)
}
