package runtime

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoopWillStopNormally(t *testing.T) {
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
			if data.value == 10 {
				cancel()
			} else {
				data.value += 1
				value += 1
			}
			return nil
		},
	}

	fnRuntime := NewLoop(loop)

	startErr := fnRuntime.Start()
	fnRuntime.(*Runner).wait.Wait()

	assert.NoError(startErr)
	assert.Equal(10, value)
	assert.Equal(1, initCalled)
	assert.Equal(0, exitCalled)
}

func TestLoopWillStopOnError(t *testing.T) {
	assert := assert.New(t)

	value := 0
	initCalled := 0
	exitCalled := 0

	stopErr := errors.New("stopping here")

	loop := &TestLoop{
		initialise: func() (*TestData, error) {
			initCalled += 1
			return &TestData{}, nil
		},
		cleanup: func(data *TestData) {
			exitCalled += 1
		},
		loop: func(data *TestData, cancel context.CancelFunc) error {
			if data.value == 11 {
				return stopErr
			} else {
				data.value += 1
				value += 1
			}
			return nil
		},
	}

	fnRuntime := NewLoop(loop)

	startErr := fnRuntime.Start()
	fnRuntime.(*Runner).wait.Wait()
	fnRuntime.Stop()

	assert.NoError(startErr)
	assert.Equal(11, value)
	assert.Equal(1, initCalled)
	assert.Equal(1, exitCalled)
}

func TestLoopWillStopOnStop(t *testing.T) {
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

	startErr := fnRuntime.Start()
	time.Sleep(time.Millisecond)
	fnRuntime.Stop()
	fnRuntime.Stop() // testing multiple stop will not break the system

	assert.NoError(startErr)
	assert.Greater(value, 0)
	assert.Equal(1, initCalled)
	assert.Equal(1, exitCalled)
}

func TestLoopMissingLoop(t *testing.T) {
	assert := assert.New(t)

	fnRuntime := NewLoop(nil)

	startErr := fnRuntime.Start()
	fnRuntime.Stop()
	assert.ErrorIs(startErr, ErrLoopRuntimeNoLoop)
}

func TestLoopInitialiseError(t *testing.T) {
	assert := assert.New(t)

	value := 0
	initCalled := 0
	exitCalled := 0

	errorInit := errors.New("error init")

	loop := &TestLoop{
		initialise: func() (*TestData, error) {
			initCalled += 1
			return &TestData{}, errorInit
		},
		cleanup: func(data *TestData) {
			exitCalled += 1
		},
		loop: func(data *TestData, cancel context.CancelFunc) error {
			if data.value == 10 {
				cancel()
			} else {
				data.value += 1
				value += 1
			}
			return nil
		},
	}
	fnRuntime := NewLoop(loop)

	startErr := fnRuntime.Start()
	fnRuntime.Stop()

	assert.ErrorIs(startErr, ErrLoopRuntimeInitialise)
	assert.ErrorIs(startErr, errorInit)
}

// test helpers

type TestData struct {
	value int64
}

type TestLoop struct {
	data       *TestData
	initialise func() (*TestData, error)
	cleanup    func(*TestData)
	loop       func(*TestData, context.CancelFunc) error
}

func (l *TestLoop) Start() error {
	data, initErr := l.initialise()
	l.data = data
	return initErr
}
func (l *TestLoop) Stop() {
	l.cleanup(l.data)
}
func (l *TestLoop) Loop(ctx context.Context, cancel context.CancelFunc) error {
	return l.loop(l.data, cancel)
}
