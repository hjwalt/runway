package inverse

import (
	"context"
	"errors"

	"github.com/hjwalt/runway/reflect"
	"github.com/hjwalt/runway/runtime"
)

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
