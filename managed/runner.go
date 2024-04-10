package managed

import (
	"context"
	"errors"
	"sync"

	"github.com/hjwalt/runway/inverse"
)

func NewRunner(runnable Runnable) Service {
	return &runner{
		runnable: runnable,
	}
}

type runner struct {
	wait      sync.WaitGroup
	lifecycle Lifecycle
	runnable  Runnable
}

func (c *runner) Name() string {
	return "runner"
}

func (r *runner) Register(ctx context.Context, ic inverse.Container) error {
	return nil
}

func (r *runner) Resolve(ctx context.Context, ic inverse.Container) error {
	lifecycle, lifecycleErr := GetLifecycle(ic, ctx)
	if lifecycleErr != nil {
		return lifecycleErr
	}
	r.lifecycle = lifecycle

	return nil
}

func (r *runner) Clean() error {
	return nil
}

func (r *runner) Start() error {
	if r.runnable == nil {
		return ErrRunnerRuntimeNoRunnable
	}
	go r.Run()
	r.wait.Add(1)

	return nil
}

func (r *runner) Stop() error {
	r.wait.Wait()
	return nil
}

func (r *runner) Run() {
	defer r.wait.Done()
	if err := r.runnable.Run(); err != nil {
		r.lifecycle.Error(err)
	}
}

// Errors
var (
	ErrRunnerRuntimeNoRunnable    = errors.New("runner runtime no runnable provided")
	ErrRunnerRuntimeRunnableStart = errors.New("runner runtime runnable failed to start")
)
