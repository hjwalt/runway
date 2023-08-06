package inverse

import (
	"context"
	"sync"
)

type Injector[T any] func(context.Context) (T, error)

var qualifierMutex = sync.Mutex{}

var qualifierInjectorMap = map[string][]Injector[any]{}

var qualifierInjectedLast = map[string]any{}

var qualifierInjectedAll = map[string][]any{}

func Reset() {
	qualifierMutex.Lock()
	defer qualifierMutex.Unlock()

	qualifierInjectorMap = map[string][]Injector[any]{}
	Release()
}

func Release() {
	qualifierInjectedLast = map[string]any{}
	qualifierInjectedAll = map[string][]any{}
}

func Register[T any](qualifier string, injector Injector[T]) {
	qualifierMutex.Lock()
	defer qualifierMutex.Unlock()

	qualifierList, qualifierExist := qualifierInjectorMap[qualifier]
	if !qualifierExist {
		qualifierList = make([]Injector[any], 0)
	}

	qualifierList = append(qualifierList, AnyInjector(injector))
	qualifierInjectorMap[qualifier] = qualifierList
}

func AnyInjector[T any](injector Injector[T]) Injector[any] {
	return func(ctx context.Context) (any, error) {
		return injector(ctx)
	}
}
