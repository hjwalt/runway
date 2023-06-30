package runtime_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/hjwalt/runway/reflect"
	"github.com/hjwalt/runway/runtime"
	"github.com/stretchr/testify/assert"
)

type TestContextKey struct {
	key string
}

func ContextKey(key string) TestContextKey {
	return TestContextKey{key: key}
}

type TestData struct {
	value int64
}

func NewController() runtime.Controller {
	controller, _ := runtime.NewPrimaryController()
	return controller
}

func TestFunctionalWillStopNormally(t *testing.T) {
	assert := assert.New(t)

	controller := NewController()
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
	fnRuntime.SetController(controller)

	startErr := fnRuntime.Start()
	controller.Wait()

	assert.NoError(startErr)
	assert.Equal(10, value)
	assert.Equal(1, initCalled)
	assert.Equal(1, exitCalled)
}

func TestFunctionalWithContext(t *testing.T) {
	assert := assert.New(t)

	controller := NewController()
	value := int64(0)
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
			increment := reflect.GetInt64(ctx.Value(ContextKey("test")))
			data.value += increment
			value += increment
			cancel()
			return nil
		}),
		runtime.FunctionalWithContext[*TestData](context.WithValue(context.Background(), ContextKey("test"), int64(5))),
	)
	fnRuntime.SetController(controller)

	startErr := fnRuntime.Start()
	controller.Wait()

	assert.NoError(startErr)
	assert.Equal(int64(5), value)
	assert.Equal(1, initCalled)
	assert.Equal(1, exitCalled)
}

func TestFunctionalMissingInitialiseWillConstruct(t *testing.T) {
	assert := assert.New(t)

	controller := NewController()
	value := 0
	exitCalled := 0

	fnRuntime := runtime.NewFunctional[*TestData](
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
	fnRuntime.SetController(controller)

	startErr := fnRuntime.Start()
	controller.Wait()

	assert.NoError(startErr)
	assert.Equal(10, value)
	assert.Equal(1, exitCalled)
}

func TestFunctionalMissingCleanup(t *testing.T) {
	assert := assert.New(t)

	controller := NewController()
	value := 0
	initCalled := 0

	fnRuntime := runtime.NewFunctional[*TestData](
		runtime.FunctionalWithInitialise[*TestData](func() (*TestData, error) {
			initCalled += 1
			return &TestData{}, nil
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
	fnRuntime.SetController(controller)

	startErr := fnRuntime.Start()
	controller.Wait()

	assert.NoError(startErr)
	assert.Equal(10, value)
	assert.Equal(1, initCalled)
}

func TestFunctionalWillStopOnError(t *testing.T) {
	assert := assert.New(t)

	controller := NewController()
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
				return errors.New("test error")
			} else {
				data.value += 1
				value += 1
			}
			return nil
		}),
	)
	fnRuntime.SetController(controller)

	startErr := fnRuntime.Start()
	controller.Wait()

	assert.NoError(startErr)
	assert.Equal(10, value)
	assert.Equal(1, initCalled)
	assert.Equal(1, exitCalled)
}

func TestFunctionalWillStopOnStop(t *testing.T) {
	assert := assert.New(t)

	controller := NewController()
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
	fnRuntime.SetController(controller)

	startErr := fnRuntime.Start()
	time.Sleep(time.Millisecond)
	fnRuntime.Stop()
	controller.Wait()

	assert.NoError(startErr)
	assert.Greater(value, 0)
	assert.Equal(1, initCalled)
	assert.Equal(1, exitCalled)
}

func TestFunctionalMissingLoop(t *testing.T) {
	assert := assert.New(t)

	controller := NewController()
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
	)
	fnRuntime.SetController(controller)

	startErr := fnRuntime.Start()
	assert.ErrorIs(startErr, runtime.ErrFunctionalRuntimeNoLoop)
}

func TestFunctionalInitialiseError(t *testing.T) {
	assert := assert.New(t)

	controller := NewController()
	value := 0
	initCalled := 0
	exitCalled := 0

	fnRuntime := runtime.NewFunctional[*TestData](
		runtime.FunctionalWithInitialise[*TestData](func() (*TestData, error) {
			initCalled += 1
			return &TestData{}, errors.New("error init")
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
	fnRuntime.SetController(controller)

	startErr := fnRuntime.Start()
	assert.ErrorIs(startErr, runtime.ErrFunctionalRuntimeInitialise)
}
