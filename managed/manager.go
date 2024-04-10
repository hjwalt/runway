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
	ctx context.Context,
	ic inverse.Container,
	services []Service,
	components []Component,
	configurations []Configuration,
) runtime.Runtime {
	return &manager{
		ctx: ctx,
		ic:  ic,
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
	ctx       context.Context
	ic        inverse.Container
	lifecycle Lifecycle
}

func (r *manager) Start() error {
	if err := r.lifecycle.Register(r.ctx, r.ic); err != nil {
		return err
	}
	if err := r.lifecycle.Resolve(r.ctx, r.ic); err != nil {
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
