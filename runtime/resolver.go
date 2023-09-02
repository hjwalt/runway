package runtime

import (
	"context"
	"errors"

	"github.com/hjwalt/runway/inverse"
	"github.com/hjwalt/runway/reflect"
)

type Resolver[T any] interface {
	AddConfig(inverse.Injector[Configuration[T]])
	AddConfigVal(Configuration[T])
	Register()
}

func NewResolver[T any, I any](
	qualifier string,
	container inverse.Container,
	configurationRequired bool,
	constructor Constructor[T, I],
) Resolver[T] {
	return &resolver[T, I]{
		qualifier:             qualifier,
		container:             container,
		configurationRequired: configurationRequired,
		constructor:           constructor,
	}
}

type resolver[T any, I any] struct {
	qualifier             string
	container             inverse.Container
	configurationRequired bool
	constructor           Constructor[T, I]
}

func (r *resolver[T, I]) AddConfig(injector inverse.Injector[Configuration[T]]) {
	inverse.GenericAdd[Configuration[T]](r.container, QualifierConfig(r.qualifier), injector)
}

func (r *resolver[T, I]) AddConfigVal(config Configuration[T]) {
	r.container.AddVal(QualifierConfig(r.qualifier), config)
}

func (r *resolver[T, I]) Register() {
	r.container.Add(r.qualifier, func(ctx context.Context, ci inverse.Container) (any, error) {
		configurations, getConfigurationError := inverse.GenericGetAll[Configuration[T]](ci, ctx, QualifierConfig(r.qualifier))
		if getConfigurationError != nil {
			if r.configurationRequired || !errors.Is(getConfigurationError, inverse.ErrInverseResolverMissing) {
				return reflect.Construct[T](), getConfigurationError
			}
		}
		return r.constructor(configurations...), nil
	})
}

func QualifierConfig(q string) string {
	return q + "Config"
}
