package managed_test

import (
	"context"
	"errors"
	"testing"

	"github.com/hjwalt/runway/inverse"
	"github.com/hjwalt/runway/managed"
	"github.com/stretchr/testify/assert"
)

func TestChannelWillStopNormally(t *testing.T) {
	assert := assert.New(t)

	value := 0

	loop := &TestChannel{
		loop: func(data *TestData) error {
			data.value += 1
			value += 1
			return nil
		},
	}

	fnRuntime := managed.NewChannel[*TestData](loop)

	services := []managed.Service{}
	services = append(services, fnRuntime)
	services = append(services, loop)
	manager := managed.New(services, []managed.Component{}, []managed.Configuration{})

	startErr := manager.Start()
	assert.NoError(startErr)

	loop.channel <- &TestData{}
	loop.channel <- &TestData{}
	loop.channel <- &TestData{}
	loop.channel <- &TestData{}
	loop.channel <- &TestData{}
	loop.channel <- &TestData{}
	loop.channel <- &TestData{}
	loop.channel <- &TestData{}
	loop.channel <- &TestData{}
	loop.channel <- &TestData{}

	manager.Stop()

	assert.Equal(10, value)
}

func TestChannelBrokenInit(t *testing.T) {
	assert := assert.New(t)

	value := 0

	loop := &TestChannel{
		brokenInit: true,
		loop: func(data *TestData) error {
			data.value += 1
			value += 1
			return nil
		},
	}

	fnRuntime := managed.NewChannel[*TestData](loop)

	services := []managed.Service{}
	services = append(services, fnRuntime)
	services = append(services, loop)
	manager := managed.New(services, []managed.Component{}, []managed.Configuration{})

	startErr := manager.Start()
	assert.Error(startErr)

	manager.Stop()
	assert.Equal(0, value)
}

func TestChannelBrokenChannel(t *testing.T) {
	assert := assert.New(t)

	value := 0

	loop := &TestChannel{
		brokenChannel: true,
		loop: func(data *TestData) error {
			data.value += 1
			value += 1
			return nil
		},
	}

	fnRuntime := managed.NewChannel[*TestData](loop)

	services := []managed.Service{}
	services = append(services, fnRuntime)
	services = append(services, loop)
	manager := managed.New(services, []managed.Component{}, []managed.Configuration{})

	startErr := manager.Start()
	assert.NoError(startErr)
	manager.Stop()
	assert.Equal(0, value)
}

func TestChannelNilChannel(t *testing.T) {
	assert := assert.New(t)

	value := 0

	fnRuntime := managed.NewChannel[*TestData](nil)

	services := []managed.Service{}
	services = append(services, fnRuntime)
	manager := managed.New(services, []managed.Component{}, []managed.Configuration{})

	startErr := manager.Start()
	assert.NoError(startErr)
	manager.Stop()
	assert.Equal(0, value)
}

func TestChannelBrokenLoop(t *testing.T) {
	assert := assert.New(t)

	value := 0

	loop := &TestChannel{
		loop: func(data *TestData) error {
			data.value += 1
			value += 1
			if value == 3 {
				return errors.New("broken value")
			}
			return nil
		},
	}

	fnRuntime := managed.NewChannel[*TestData](loop)

	services := []managed.Service{}
	services = append(services, fnRuntime)
	services = append(services, loop)
	manager := managed.New(services, []managed.Component{}, []managed.Configuration{})

	startErr := manager.Start()

	loop.channel <- &TestData{}
	loop.channel <- &TestData{}
	loop.channel <- &TestData{}

	manager.Stop()

	assert.NoError(startErr)
	assert.Equal(3, value)
}

type TestChannel struct {
	brokenInit    bool
	brokenChannel bool
	channel       chan *TestData
	loop          func(*TestData) error
}

func (r *TestChannel) Name() string {
	return "test-channel"
}

func (r *TestChannel) Register(ctx context.Context, ic inverse.Container) error {
	return nil
}

func (r *TestChannel) Resolve(ctx context.Context, ic inverse.Container) error {
	return nil
}

func (r *TestChannel) Clean() error {
	return nil
}

func (l *TestChannel) Start() error {
	l.channel = make(chan *TestData)
	if l.brokenInit {
		return errors.New("broken init")
	}
	return nil
}

func (l *TestChannel) Stop() error {
	alreadyClosed := false
	defer func() {
		if recover() != nil {
			// The return result can be altered
			// in a defer function call.
			alreadyClosed = true
		}
	}()

	close(l.channel)

	if alreadyClosed {
		return errors.New("channel already closed")
	} else {
		return nil
	}
}

func (l *TestChannel) Channel() (<-chan *TestData, error) {
	if l.brokenChannel {
		return nil, errors.New("broken channel")
	}
	return l.channel, nil
}

func (l *TestChannel) Loop(d *TestData) error {
	return l.loop(d)
}

type TestData struct {
	value int
}
