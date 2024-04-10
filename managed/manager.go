package managed

import (
	"context"

	"github.com/hjwalt/runway/inverse"
)

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
