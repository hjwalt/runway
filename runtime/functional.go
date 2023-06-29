package runtime

import (
	"context"
	"errors"

	"github.com/hjwalt/runway/logger"
	"github.com/hjwalt/runway/reflect"
)

// constructor
func NewFunctional[C any](configurations ...Configuration[*Functional[C]]) Runtime {
	c := &Functional[C]{}
	c = FunctionalDefault[C](c)
	for _, configuration := range configurations {
		c = configuration(c)
	}
	return c
}

// default
func FunctionalDefault[C any](c *Functional[C]) *Functional[C] {
	c.data = reflect.Construct[C]()

	ctx, cancel := context.WithCancel(context.Background())
	c.context = ctx
	c.cancel = cancel

	return c
}

// configuration
func FunctionalWithInitialise[C any](initialise func() (C, error)) Configuration[*Functional[C]] {
	return func(c *Functional[C]) *Functional[C] {
		c.initialise = initialise
		return c
	}
}

func FunctionalWithCleanup[C any](cleanup func(C)) Configuration[*Functional[C]] {
	return func(c *Functional[C]) *Functional[C] {
		c.cleanup = cleanup
		return c
	}
}

func FunctionalWithLoop[C any](loop func(data C, ctx context.Context, cancel context.CancelFunc) error) Configuration[*Functional[C]] {
	return func(c *Functional[C]) *Functional[C] {
		c.loop = loop
		return c
	}
}

func FunctionalWithContext[C any](ctx context.Context) Configuration[*Functional[C]] {
	return func(c *Functional[C]) *Functional[C] {
		ctx, cancel := context.WithCancel(ctx)
		c.context = ctx
		c.cancel = cancel
		return c
	}
}

// implementation
type Functional[C any] struct {
	controller Controller
	initialise func() (C, error)
	cleanup    func(C)
	loop       func(data C, ctx context.Context, cancel context.CancelFunc) error
	context    context.Context
	cancel     context.CancelFunc
	data       C
}

func (r *Functional[C]) Start() error {
	if r.loop == nil {
		return ErrFunctionalRuntimeNoLoop
	}

	if r.initialise != nil {
		data, initerr := r.initialise()
		if initerr != nil {
			return errors.Join(ErrFunctionalRuntimeInitialise, initerr)
		}
		r.data = data
	}

	go r.Run()
	r.controller.Started()
	return nil
}

func (r *Functional[C]) Stop() {
	r.cancel()
}

func (r *Functional[C]) SetController(controller Controller) {
	r.controller = controller
}

func (r *Functional[C]) Run() {
	defer r.controller.Stopped()

	for {
		err := r.loop(r.data, r.context, r.cancel)
		if err != nil {
			logger.ErrorErr("functional runtime loop error", err)
			r.controller.Error(err)
			break
		}
		if r.context.Err() != nil {
			logger.WarnErr("functional runtime exitting via context", err)
			break
		}
	}

	r.cleanup(r.data)
}

// Errors
var (
	ErrFunctionalRuntimeNoLoop     = errors.New("functional runtime no loop function provided")
	ErrFunctionalRuntimeInitialise = errors.New("functional runtime initialise function failed")
)
