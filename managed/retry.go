package managed

import (
	"context"
	"errors"
	"time"

	retrylib "github.com/avast/retry-go/v4"
	"github.com/hjwalt/runway/inverse"
)

const (
	ConfRetryAttempts             = "ConfRetryAttempts"
	ConfRetryDelayMillisecond     = "ConfRetryDelayMillisecond"
	ConfRetryMaxDelayMillisecond  = "ConfRetryMaxDelayMillisecond"
	ConfRetryMaxJitterMillisecond = "ConfRetryMaxJitterMillisecond"
	ConfRetryAbsorbError          = "ConfRetryAbsorbError"
)

func NewRetry() Retry {
	return NewQualifiedRetry(QualifierRetry)
}

func NewQualifiedRetry(qualifier string) Retry {
	return &retry{
		qualifier: qualifier,
		options:   []retrylib.Option{},
		absorb:    false,
	}
}

func ResolveRetry(ctx context.Context, container inverse.Container, qualifier string) (Retry, error) {
	return inverse.GenericGet[Retry](container, ctx, qualifier)
}

// implementation
type retry struct {
	qualifier string
	lifecycle Lifecycle
	options   []retrylib.Option
	absorb    bool
}

func (c *retry) Name() string {
	return c.qualifier
}

func (r *retry) Register(ctx context.Context, ic inverse.Container) error {
	ic.AddVal(r.Name(), r)
	return nil
}

func (r *retry) Resolve(ctx context.Context, ic inverse.Container) error {
	config, configErr := ResolveConfig(ctx, ic, r.Name())
	if configErr != nil {
		return configErr
	}

	r.options = []retrylib.Option{
		retrylib.RetryIf(AlwaysTry),
		retrylib.DelayType(retrylib.BackOffDelay),
		retrylib.Attempts(uint(config.GetUint64(ConfRetryAttempts, 1000000))),
		retrylib.Delay(time.Duration(config.GetInt64(ConfRetryDelayMillisecond, 10)) * time.Millisecond),
		retrylib.MaxDelay(time.Duration(config.GetInt64(ConfRetryMaxDelayMillisecond, 1000)) * time.Millisecond),
		retrylib.MaxJitter(time.Duration(config.GetInt64(ConfRetryMaxJitterMillisecond, 1000)) * time.Millisecond),
	}
	r.absorb = config.GetBool(ConfRetryAbsorbError, false)

	lifecycle, lifecycleErr := ResolveLifecycle(ctx, ic)
	if lifecycleErr != nil {
		return lifecycleErr
	}
	r.lifecycle = lifecycle

	return nil
}

func (r *retry) Clean() error {
	return nil
}

func (r *retry) Do(fnToDo func(int64) error) error {
	tryCount := int64(0)

	err := retrylib.Do(func() error {
		if !r.lifecycle.Running() {
			return ErrRetryStopped
		}
		tryCount += 1
		return fnToDo(tryCount)
	}, r.options...)

	if r.absorb && err != nil {
		return nil
	}

	return err
}

// utils

var ErrRetryStopped = errors.New("retry stopped")

func AlwaysTry(err error) bool {
	if err != nil && errors.Is(err, ErrRetryStopped) {
		return false
	}
	return true
}
