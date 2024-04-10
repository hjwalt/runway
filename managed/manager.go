package managed

import (
	"context"
	"sync/atomic"

	"github.com/hjwalt/runway/inverse"
	"github.com/hjwalt/runway/runtime"
)

const (
	ConfManagerName             = "ConfManagerName"
	ConfManagerRuntimeQualifier = "ConfManagerRuntimeQualifier"
	ConfManagerQualifier        = "ConfManagerQualifier"

	DefaultManagerRuntimeQualifier = "ManagedRuntime"
)

func New(
	services []Service,
	components []Component,
	configurations []Configuration,
) runtime.Runtime {
	return &manager{
		lifecycle: &lifecycle{
			running:        &atomic.Bool{},
			services:       services,
			components:     components,
			configurations: configurations,
		},
	}
}

// implementation
type manager struct {
	lifecycle Lifecycle
}

func (r *manager) Start() error {
	ctx := context.Background()
	ic := inverse.NewContainer()

	if err := r.lifecycle.Register(ctx, ic); err != nil {
		return err
	}
	if err := r.lifecycle.Resolve(ctx, ic); err != nil {
		return err
	}
	return r.lifecycle.Start()
}

func (r *manager) Stop() {
	r.lifecycle.Clean()
	if err := r.lifecycle.Stop(); err != nil {
		panic(err)
	}
}

// Errors
