package runtime

import (
	"context"
	"errors"
	"syscall"
	"testing"
	"time"

	"github.com/hjwalt/runway/trusted"
	"github.com/stretchr/testify/assert"
)

func TestPrimaryWillStopNormally(t *testing.T) {
	assert := assert.New(t)

	value := 0
	initCalled := 0
	exitCalled := 0

	loop := &TestLoop{
		initialise: func() (*TestData, error) {
			initCalled += 1
			return &TestData{}, nil
		},
		cleanup: func(data *TestData) {
			exitCalled += 1
		},
		loop: func(data *TestData, cancel context.CancelFunc) error {
			data.value += 1
			value += 1
			return nil
		},
	}
	fnRuntime := NewLoop(loop)

	startErr := Start([]Runtime{
		fnRuntime,
	}, time.Millisecond)
	time.Sleep(time.Millisecond)
	go Stop()

	Wait()
	Stop() // testing multiple stop will not break the system

	assert.NoError(startErr)
	assert.Greater(value, 0)
	assert.Equal(1, initCalled)
	assert.Equal(1, exitCalled)
}

func TestPrimaryInitialiseErrorWillNotStart(t *testing.T) {
	assert := assert.New(t)

	value := 0
	initCalled := 0
	exitCalled := 0

	initErr := errors.New("bulbasaur")

	loop := &TestLoop{
		initialise: func() (*TestData, error) {
			initCalled += 1
			return &TestData{}, initErr
		},
		cleanup: func(data *TestData) {
			exitCalled += 1
		},
		loop: func(data *TestData, cancel context.CancelFunc) error {
			data.value += 1
			value += 1
			return nil
		},
	}
	fnRuntime := NewLoop(loop)

	startErr := Start([]Runtime{
		fnRuntime,
	}, time.Millisecond)
	Wait()

	assert.ErrorIs(startErr, initErr)
	assert.Equal(0, value)
	assert.Equal(1, initCalled)
	assert.Equal(0, exitCalled)
}

func TestPrimaryWillStopOnError(t *testing.T) {
	assert := assert.New(t)

	value := 0
	initCalled := 0
	exitCalled := 0

	loop := &TestLoop{
		initialise: func() (*TestData, error) {
			initCalled += 1
			return &TestData{}, nil
		},
		cleanup: func(data *TestData) {
			exitCalled += 1
		},
		loop: func(data *TestData, cancel context.CancelFunc) error {
			data.value += 1
			value += 1
			return trusted.ErrPrimaryTesting
		},
	}
	fnRuntime := NewLoop(loop)

	startErr := Start([]Runtime{
		fnRuntime,
	}, time.Millisecond)
	Wait()

	assert.NoError(startErr)
	assert.Equal(1, value)
	assert.Equal(1, initCalled)
	assert.Equal(1, exitCalled)
}

func TestPrimaryWillStopOnSignal(t *testing.T) {
	assert := assert.New(t)

	value := 0
	initCalled := 0
	exitCalled := 0

	loop := &TestLoop{
		initialise: func() (*TestData, error) {
			initCalled += 1
			return &TestData{}, nil
		},
		cleanup: func(data *TestData) {
			exitCalled += 1
		},
		loop: func(data *TestData, cancel context.CancelFunc) error {
			data.value += 1
			value += 1
			return nil
		},
	}
	fnRuntime := NewLoop(loop)

	startErr := Start([]Runtime{
		fnRuntime,
	}, time.Millisecond)
	time.Sleep(time.Millisecond)
	primary.Load().interruptChannel <- syscall.SIGTERM
	Wait()

	assert.NoError(startErr)
	assert.Greater(value, 0)
	assert.Equal(1, initCalled)
	assert.Equal(1, exitCalled)
}
