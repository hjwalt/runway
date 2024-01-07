package runtime

import (
	"context"
	"errors"
	"sync/atomic"
	"time"

	retry "github.com/avast/retry-go/v4"
	"github.com/hjwalt/runway/logger"
	"github.com/hjwalt/runway/reflect"
)

// constructor
func NewRetry(configurations ...Configuration[*Retry]) *Retry {
	consumer := &Retry{
		options: []retry.Option{
			retry.RetryIf(AlwaysTry),
			retry.Attempts(1000000), // more than a week as a default which is practically infinite
			retry.Delay(10 * time.Millisecond),
			retry.MaxDelay(time.Second),
			retry.MaxJitter(time.Second),
			retry.DelayType(retry.BackOffDelay),
		},
		absorb: false,
	}
	for _, configuration := range configurations {
		consumer = configuration(consumer)
	}
	return consumer
}

// configurations
func WithRetryAttempts(attempts uint) Configuration[*Retry] {
	return func(c *Retry) *Retry {
		c.options = append(c.options, retry.Attempts(attempts))
		return c
	}
}

func WithRetryDelay(delay time.Duration) Configuration[*Retry] {
	return func(c *Retry) *Retry {
		c.options = append(c.options, retry.Delay(delay))
		return c
	}
}

func WithRetryMaxDelay(delay time.Duration) Configuration[*Retry] {
	return func(c *Retry) *Retry {
		c.options = append(c.options, retry.MaxDelay(delay))
		return c
	}
}

func WithRetryAbsorbError(absorb bool) Configuration[*Retry] {
	return func(c *Retry) *Retry {
		c.absorb = absorb
		return c
	}
}

// implementation
type Retry struct {
	options []retry.Option
	absorb  bool
	stopped atomic.Bool
}

func (c *Retry) Start() error {
	c.stopped.Store(false)
	logger.Debug("retry start")
	return nil
}

func (c *Retry) Stop() {
	c.stopped.Store(true)
	logger.Debug("retry stop")
}

func (c *Retry) Do(fnToDo func(int64) error) error {
	tryCount := int64(0)

	err := retry.Do(func() error {
		if c.stopped.Load() {
			return ErrRetryStopped
		}
		tryCount += 1
		return fnToDo(tryCount)
	}, c.options...)

	if c.absorb && err != nil {
		logger.ErrorErr("absorbing retry error", err)
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

func SetRetryCount(ctx context.Context, trycount int64) context.Context {
	return context.WithValue(ctx, Context("RetryCount"), trycount)
}

func GetRetryCount(ctx context.Context) int64 {
	return reflect.GetInt64(ctx.Value(Context("RetryCount")))
}
