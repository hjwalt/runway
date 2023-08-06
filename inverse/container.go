package inverse

import (
	"context"
	"errors"
	"sync"

	"github.com/hjwalt/runway/reflect"
	"github.com/hjwalt/runway/runtime"
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

func RegisterInstance[T any](qualifier string, instance T) {
	Register[T](qualifier, func(ctx context.Context) (T, error) { return instance, nil })
}

func RegisterInstances[T any](qualifier string, instances []T) {
	for _, instance := range instances {
		RegisterInstance[T](qualifier, instance)
	}
}

func AnyInjector[T any](injector Injector[T]) Injector[any] {
	return func(ctx context.Context) (any, error) {
		return injector(ctx)
	}
}

func RegisterConfiguration[T any](configurationQualifier string, configuration runtime.Configuration[T]) {
	RegisterInstance[runtime.Configuration[T]](configurationQualifier, configuration)
}

func RegisterConfigurations[T any](configurationQualifier string, configurations []runtime.Configuration[T]) {
	RegisterInstances[runtime.Configuration[T]](configurationQualifier, configurations)
}

func RegisterWithConfigurationOptional[T any, I any](qualifier string, configurationQualifier string, constructor runtime.Constructor[T, I]) {
	Register[I](qualifier, func(ctx context.Context) (I, error) {
		configurations, getConfigurationError := GetAll[runtime.Configuration[T]](ctx, configurationQualifier)
		if getConfigurationError != nil && !errors.Is(getConfigurationError, ErrNotInjected) {
			return reflect.Construct[I](), getConfigurationError
		}
		return constructor(configurations...), nil
	})
}

func RegisterWithConfigurationRequired[T any, I any](qualifier string, configurationQualifier string, constructor runtime.Constructor[T, I]) {
	Register[I](qualifier, func(ctx context.Context) (I, error) {
		configurations, getConfigurationError := GetAll[runtime.Configuration[T]](ctx, configurationQualifier)
		if getConfigurationError != nil {
			return reflect.Construct[I](), getConfigurationError
		}
		return constructor(configurations...), nil
	})
}
