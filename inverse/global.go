package inverse

import (
	"context"
)

type Injector[T any] func(context.Context, Container) (T, error)

var Global = NewContainer()

func Reset() {
	Global.Reset()
}

func Release() {
	Global.Clear()
}

func Register[T any](qualifier string, injector Injector[T]) {
	GenericAdd(Global, qualifier, injector)
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
