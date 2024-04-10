package managed

import (
	"context"
	"errors"

	"github.com/hjwalt/runway/inverse"
)

func NewLoop(loop Loop) []Service {
	c := &loopRunnable{
		loop: loop,
	}
	return []Service{c, NewRunner(c)}
}

type loopRunnable struct {
	loop    Loop
	context context.Context
	cancel  context.CancelFunc
}

func (c *loopRunnable) Name() string {
	return "loop"
}

func (r *loopRunnable) Register(ctx context.Context, ic inverse.Container) error {
	if r.loop == nil {
		return ErrLoopRuntimeNoLoop
	}

	ctx, cancel := context.WithCancel(ctx)
	r.context = ctx
	r.cancel = cancel

	return nil
}

func (r *loopRunnable) Resolve(ctx context.Context, ic inverse.Container) error {
	return nil
}

func (r *loopRunnable) Clean() error {
	return nil
}

func (r *loopRunnable) Start() error {
	return nil
}

func (r *loopRunnable) Stop() error {
	r.cancel()
	return nil
}

func (r *loopRunnable) Run() error {
	for {
		err := r.loop.Loop()
		if err != nil {
			return err
		}
		if r.context.Err() != nil {
			return nil
		}
	}
}

// Errors
var (
	ErrLoopRuntimeNoLoop     = errors.New("functional runtime no loop function provided")
	ErrLoopRuntimeInitialise = errors.New("functional runtime initialise function failed")
)
