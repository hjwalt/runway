package runtime

import (
	"context"
	"errors"

	"github.com/hjwalt/runway/logger"
)

type Loop interface {
	Start() error
	Stop()
	Loop(context.Context, context.CancelFunc) error
}

// constructor
func NewLoop(loop Loop) Runtime {
	c := &LoopRunnable{
		loop: loop,
	}

	ctx, cancel := context.WithCancel(context.Background())
	c.context = ctx
	c.cancel = cancel

	return NewRunner(c)
}

// implementation
type LoopRunnable struct {
	loop    Loop
	context context.Context
	cancel  context.CancelFunc
}

func (r *LoopRunnable) Start() error {
	if r.loop == nil {
		return ErrLoopRuntimeNoLoop
	}

	initerr := r.loop.Start()
	if initerr != nil {
		return errors.Join(ErrLoopRuntimeInitialise, initerr)
	}

	return nil
}

func (r *LoopRunnable) Stop() {
	if r.loop != nil {
		r.loop.Stop()
	}
	r.cancel()
}

func (r *LoopRunnable) Run() error {
	for {
		err := r.loop.Loop(r.context, r.cancel)
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
