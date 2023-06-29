package runtime_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/hjwalt/runway/runtime"
	"github.com/stretchr/testify/assert"
)

func TestPrimaryWillStopNormally(t *testing.T) {
	assert := assert.New(t)

	value := 0
	initCalled := 0
	exitCalled := 0

	fnRuntime := runtime.NewFunctional[*TestData](
		runtime.FunctionalWithInitialise[*TestData](func() (*TestData, error) {
			initCalled += 1
			return &TestData{}, nil
		}),
		runtime.FunctionalWithCleanup[*TestData](func(data *TestData) {
			exitCalled += 1
		}),
		runtime.FunctionalWithLoop[*TestData](func(data *TestData, ctx context.Context, cancel context.CancelFunc) error {
			if data.value == 10 {
				cancel()
			} else {
				data.value += 1
				value += 1
			}
			return nil
		}),
	)

	fnPrimary := runtime.NewPrimary(
		runtime.PrimaryWithRuntime(fnRuntime),
	)

	startErr := fnPrimary.Start()

	assert.NoError(startErr)
	assert.Equal(10, value)
	assert.Equal(1, initCalled)
	assert.Equal(1, exitCalled)
}

func TestPrimaryWillNotStartOnOneInitError(t *testing.T) {
	assert := assert.New(t)

	value := 0
	initCalled := 0
	exitCalled := 0

	fnRuntime1 := runtime.NewFunctional[*TestData](
		runtime.FunctionalWithInitialise[*TestData](func() (*TestData, error) {
			initCalled += 1
			return &TestData{}, nil
		}),
		runtime.FunctionalWithCleanup[*TestData](func(data *TestData) {
			exitCalled += 1
		}),
		runtime.FunctionalWithLoop[*TestData](func(data *TestData, ctx context.Context, cancel context.CancelFunc) error {
			if data.value == 10 {
				cancel()
			} else {
				data.value += 1
				value += 1
			}
			return nil
		}),
	)

	fnRuntime2 := runtime.NewFunctional[*TestData](
		runtime.FunctionalWithInitialise[*TestData](func() (*TestData, error) {
			initCalled += 1
			return &TestData{}, errors.New("test error")
		}),
		runtime.FunctionalWithCleanup[*TestData](func(data *TestData) {
			exitCalled += 1
		}),
		runtime.FunctionalWithLoop[*TestData](func(data *TestData, ctx context.Context, cancel context.CancelFunc) error {
			data.value += 1
			value += 1
			return nil
		}),
	)

	fnPrimary := runtime.NewPrimary(
		runtime.PrimaryWithRuntime(fnRuntime1),
		runtime.PrimaryWithRuntime(fnRuntime2),
	)

	startErr := fnPrimary.Start()

	assert.ErrorIs(startErr, runtime.ErrPrimaryInitialiseError)
	assert.Greater(value, 0)
	assert.Equal(2, initCalled)
	assert.Equal(1, exitCalled)
}

func TestPrimaryWillStopWhenStopped(t *testing.T) {
	assert := assert.New(t)

	value := 0
	initCalled := 0
	exitCalled := 0

	fnRuntime := runtime.NewFunctional[*TestData](
		runtime.FunctionalWithInitialise[*TestData](func() (*TestData, error) {
			initCalled += 1
			return &TestData{}, nil
		}),
		runtime.FunctionalWithCleanup[*TestData](func(data *TestData) {
			exitCalled += 1
		}),
		runtime.FunctionalWithLoop[*TestData](func(data *TestData, ctx context.Context, cancel context.CancelFunc) error {
			data.value += 1
			value += 1
			return nil
		}),
	)

	fnPrimary := runtime.NewPrimary(
		runtime.PrimaryWithRuntime(fnRuntime),
	)

	go func() {
		time.Sleep(time.Millisecond)
		fnPrimary.Stop()
	}()
	startErr := fnPrimary.Start()

	assert.NoError(startErr)
	assert.Greater(value, 0)
	assert.Equal(1, initCalled)
	assert.Equal(1, exitCalled)
}

func TestPrimaryWillStopWhenError(t *testing.T) {
	assert := assert.New(t)

	value := 0
	initCalled := 0
	exitCalled := 0

	fnRuntime := runtime.NewFunctional[*TestData](
		runtime.FunctionalWithInitialise[*TestData](func() (*TestData, error) {
			initCalled += 1
			return &TestData{}, nil
		}),
		runtime.FunctionalWithCleanup[*TestData](func(data *TestData) {
			exitCalled += 1
		}),
		runtime.FunctionalWithLoop[*TestData](func(data *TestData, ctx context.Context, cancel context.CancelFunc) error {
			data.value += 1
			value += 1
			return runtime.ErrPrimaryTesting
		}),
	)

	fnPrimary := runtime.NewPrimary(
		runtime.PrimaryWithRuntime(fnRuntime),
	)

	startErr := fnPrimary.Start()

	assert.NoError(startErr)
	assert.Greater(value, 0)
	assert.Equal(1, initCalled)
	assert.Equal(1, exitCalled)
}
