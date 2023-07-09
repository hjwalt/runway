package runtime

import (
	"context"
	"errors"

	"github.com/hjwalt/runway/logger"
	"github.com/hjwalt/runway/reflect"
)

type Loop[C any] interface {
	Initialise() (C, error)
	Cleanup(C)
	Loop(C, context.CancelFunc) error
}

// constructor
func NewLoop[C any](configurations ...Configuration[*LoopRunnable[C]]) Runtime {
	c := &LoopRunnable[C]{}
	c = LoopDefault[C](c)
	for _, configuration := range configurations {
		c = configuration(c)
	}
	return NewRunner(c)
}

// default
func LoopDefault[C any](c *LoopRunnable[C]) *LoopRunnable[C] {
	c.data = reflect.Construct[C]()

	ctx, cancel := context.WithCancel(context.Background())
	c.context = ctx
	c.cancel = cancel

	return c
}

// configuration
func LoopWithLoop[C any](loop Loop[C]) Configuration[*LoopRunnable[C]] {
	return func(c *LoopRunnable[C]) *LoopRunnable[C] {
		c.loop = loop
		return c
	}
}

// implementation
type LoopRunnable[C any] struct {
	loop    Loop[C]
	context context.Context
	cancel  context.CancelFunc
	data    C
}

func (r *LoopRunnable[C]) Start() error {
	if r.loop == nil {
		return ErrLoopRuntimeNoLoop
	}

	data, initerr := r.loop.Initialise()
	if initerr != nil {
		return errors.Join(ErrLoopRuntimeInitialise, initerr)
	}
	r.data = data

	return nil
}

func (r *LoopRunnable[C]) Stop() {
	r.cancel()
}

func (r *LoopRunnable[C]) Run() error {
	defer r.loop.Cleanup(r.data)

	for {
		err := r.loop.Loop(r.data, r.cancel)
		if err != nil {
			logger.ErrorErr("functional runtime loop error", err)
			return err
		}
		if r.context.Err() != nil {
			logger.WarnErr("functional runtime exitting via context", err)
			return nil
		}
	}
}

// Errors
var (
	ErrLoopRuntimeNoLoop     = errors.New("functional runtime no loop function provided")
	ErrLoopRuntimeInitialise = errors.New("functional runtime initialise function failed")
)
