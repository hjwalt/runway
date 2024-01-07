package runtime

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRetryNoAbsorb(t *testing.T) {
	assert := assert.New(t)

	retryRuntime := NewRetry(
		WithRetryAttempts(100),
		WithRetryDelay(1*time.Microsecond),
		WithRetryMaxDelay(1*time.Millisecond),
		WithRetryAbsorbError(false),
	)

	startErr := retryRuntime.Start()
	assert.NoError(startErr)

	iteration := 0
	lastiteration := int64(0)

	err := retryRuntime.Do(func(i int64) error {
		iteration += 1
		lastiteration = i
		return errors.New("retryable error")
	})
	retryRuntime.Stop()

	assert.Equal(iteration, 100)
	assert.Equal(lastiteration, int64(100))
	assert.ErrorContains(err, "All attempts fail")
}

func TestRetryStoppedInBetween(t *testing.T) {
	assert := assert.New(t)

	retryRuntime := NewRetry(
		WithRetryAttempts(100),
		WithRetryDelay(1*time.Millisecond),
		WithRetryMaxDelay(10*time.Millisecond),
		WithRetryAbsorbError(false),
	)

	startErr := retryRuntime.Start()
	assert.NoError(startErr)

	iteration := 0
	lastiteration := int64(0)
	retryErr := make(chan error)

	go func() {
		retryErr <- retryRuntime.Do(func(i int64) error {
			iteration += 1
			lastiteration = i
			return errors.New("retryable error")
		})
	}()

	time.Sleep(1 * time.Millisecond)
	retryRuntime.Stop()
	retryErrOut := <-retryErr

	assert.Greater(iteration, 0)
	assert.Greater(lastiteration, int64(0))
	assert.Less(iteration, 100)
	assert.Less(lastiteration, int64(100))
	assert.ErrorIs(retryErrOut, ErrRetryStopped)
}

func TestRetryAbsorb(t *testing.T) {
	assert := assert.New(t)

	retryRuntime := NewRetry(
		WithRetryAttempts(100),
		WithRetryDelay(1*time.Microsecond),
		WithRetryMaxDelay(1*time.Millisecond),
		WithRetryAbsorbError(true),
	)

	startErr := retryRuntime.Start()
	assert.NoError(startErr)

	iteration := 0
	lastiteration := int64(0)

	err := retryRuntime.Do(func(i int64) error {
		iteration += 1
		lastiteration = i
		return errors.New("retryable error")
	})
	retryRuntime.Stop()

	assert.Equal(iteration, 100)
	assert.Equal(lastiteration, int64(100))
	assert.NoError(err)
}
