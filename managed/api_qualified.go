package managed

import (
	"context"
	"errors"
	"sync/atomic"

	"github.com/hjwalt/runway/inverse"
	"github.com/hjwalt/runway/runtime"
)

const (
	QualifierConfiguration = "configuration"
	QualifierService       = "service"
	QualifierComponent     = "component"
)

func AddConfiguration(ic inverse.Container, component string, configuration map[string]string) {
	ic.AddVal(QualifierConfiguration, NewConfig(component, configuration))
}

func AddService(ic inverse.Container, svc Service) {
	ic.AddVal(QualifierService, svc)
}

func AddComponent(ic inverse.Container, comp Component) {
	ic.AddVal(QualifierComponent, comp)
}

func New(
	ic inverse.Container,
	ctx context.Context,
) (
	runtime.Runtime,
	error,
) {

	configurations, configGetErr := inverse.GenericGetAll[Configuration](ic, ctx, QualifierConfiguration)
	if configGetErr != nil && !errors.Is(configGetErr, inverse.ErrInverseResolverMissing) {
		return nil, configGetErr
	}
	components, compGetErr := inverse.GenericGetAll[Component](ic, ctx, QualifierComponent)
	if compGetErr != nil && !errors.Is(compGetErr, inverse.ErrInverseResolverMissing) {
		return nil, compGetErr
	}
	services, svcGetErr := inverse.GenericGetAll[Service](ic, ctx, QualifierService)
	if svcGetErr != nil && !errors.Is(svcGetErr, inverse.ErrInverseResolverMissing) {
		return nil, svcGetErr
	}

	if configurations == nil {
		configurations = []Configuration{}
	}
	if components == nil {
		components = []Component{}
	}
	if services == nil {
		services = []Service{}
	}

	return &manager{
		ctx: ctx,
		ic:  ic,
		lifecycle: &lifecycle{
			running:        &atomic.Bool{},
			services:       services,
			components:     components,
			configurations: configurations,
		},
	}, nil
}
