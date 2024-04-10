package managed

import (
	"context"
	"errors"
	"sync/atomic"

	"github.com/hjwalt/runway/inverse"
)

func ResolveLifecycle(ctx context.Context, container inverse.Container) (Lifecycle, error) {
	return inverse.GenericGet[Lifecycle](container, ctx, QualifierLifecycle)
}

type lifecycle struct {
	running        *atomic.Bool
	services       []Service
	components     []Component
	configurations []Configuration
}

func (c *lifecycle) Name() string {
	return QualifierLifecycle
}

func (r *lifecycle) Register(ctx context.Context, ic inverse.Container) error {
	ic.AddVal(r.Name(), r)
	for _, inst := range r.configurations {
		if err := inst.Register(ctx, ic); err != nil {
			return err
		}
	}
	for _, inst := range r.components {
		if err := inst.Register(ctx, ic); err != nil {
			return err
		}
	}
	for _, inst := range r.services {
		if err := inst.Register(ctx, ic); err != nil {
			return err
		}
	}
	return nil
}

func (r *lifecycle) Resolve(ctx context.Context, ic inverse.Container) error {
	for _, inst := range r.configurations {
		if err := inst.Resolve(ctx, ic); err != nil {
			return err
		}
	}
	for _, inst := range r.components {
		if err := inst.Resolve(ctx, ic); err != nil {
			return err
		}
	}
	for _, inst := range r.services {
		if err := inst.Resolve(ctx, ic); err != nil {
			return err
		}
	}
	return nil
}

func (r *lifecycle) Clean() error {
	cleanErrors := []error{}
	for _, inst := range r.configurations {
		if err := inst.Clean(); err != nil {
			cleanErrors = append(cleanErrors, err)
		}
	}
	for _, inst := range r.components {
		if err := inst.Clean(); err != nil {
			cleanErrors = append(cleanErrors, err)
		}
	}
	for _, inst := range r.services {
		if err := inst.Clean(); err != nil {
			cleanErrors = append(cleanErrors, err)
		}
	}
	return joinErrors(cleanErrors)
}

func (r *lifecycle) Start() error {
	for i := 0; i < len(r.services); i++ {
		if err := r.services[i].Start(); err != nil {
			errReturn := errors.Join(ErrLifecycleStartError, err)
			if stopErr := r.Stop(); stopErr != nil {
				errReturn = errors.Join(errReturn, stopErr)
			}
			return errReturn
		}
	}

	r.running.Store(true)
	return nil
}

func (r *lifecycle) Stop() error {

	stopErrors := []error{}
	for i := len(r.services); i > 0; i-- {
		if err := r.services[i-1].Stop(); err != nil {
			stopErrors = append(stopErrors, err)
		}
	}

	r.running.Store(false)
	return joinErrors(stopErrors)
}

func (c *lifecycle) Running() bool {
	return c.running.Load()
}

func (c *lifecycle) Error(error) {
	go c.Stop()
}

func joinErrors(errs []error) error {
	if len(errs) == 0 {
		return nil
	} else {
		return errors.Join(errs...)
	}
}

var (
	ErrLifecycleStartError = errors.New("lifecycle start failed")
)
