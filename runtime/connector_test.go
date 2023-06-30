package runtime_test

import (
	"errors"
	"testing"

	"github.com/hjwalt/runway/runtime"
	"github.com/stretchr/testify/assert"
)

func TestConnectorWillExitNormally(t *testing.T) {
	assert := assert.New(t)

	controller := NewController()

	initCalled := 0
	exitCalled := 0

	fnConnector := runtime.NewConnector(
		runtime.ConnectorWithInitialise[*TestData](func() (*TestData, error) {
			initCalled += 1
			return &TestData{value: 100}, nil
		}),
		runtime.ConnectorWithCleanup[*TestData](func(data *TestData) {
			exitCalled += int(data.value)
		}),
	)
	fnConnector.SetController(controller)

	startErr := fnConnector.Start()
	fnConnector.Stop()

	controller.Wait()

	assert.NoError(startErr)
	assert.Equal(1, initCalled)
	assert.Equal(100, exitCalled)
}

func TestConnectorWithoutInitialise(t *testing.T) {
	assert := assert.New(t)

	controller := NewController()

	initCalled := 0
	exitCalled := 0

	fnConnector := runtime.NewConnector(
		runtime.ConnectorWithCleanup[*TestData](func(data *TestData) {
			exitCalled += int(data.value)
		}),
	)
	fnConnector.SetController(controller)

	startErr := fnConnector.Start()
	fnConnector.Stop()

	controller.Wait()

	assert.NoError(startErr)
	assert.Equal(0, initCalled)
	assert.Equal(0, exitCalled)
}

func TestConnectorWithoutCleanup(t *testing.T) {
	assert := assert.New(t)

	controller := NewController()

	initCalled := 0

	fnConnector := runtime.NewConnector(
		runtime.ConnectorWithInitialise[*TestData](func() (*TestData, error) {
			initCalled += 1
			return &TestData{value: 100}, nil
		}),
	)
	fnConnector.SetController(controller)

	startErr := fnConnector.Start()
	fnConnector.Stop()

	controller.Wait()

	assert.NoError(startErr)
	assert.Equal(1, initCalled)
}

func TestConnectorInitErrorWillNotStart(t *testing.T) {
	assert := assert.New(t)

	controller := NewController()

	initCalled := 0
	exitCalled := 0

	fnConnector := runtime.NewConnector(
		runtime.ConnectorWithInitialise[*TestData](func() (*TestData, error) {
			initCalled += 1
			return &TestData{value: 100}, errors.New("test")
		}),
		runtime.ConnectorWithCleanup[*TestData](func(data *TestData) {
			exitCalled += int(data.value)
		}),
	)
	fnConnector.SetController(controller)

	startErr := fnConnector.Start()
	controller.Wait()

	assert.ErrorIs(startErr, runtime.ErrConnectorRuntimeInitialise)
	assert.Equal(1, initCalled)
	assert.Equal(0, exitCalled)
}
