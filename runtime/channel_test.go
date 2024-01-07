package runtime

import (
	"errors"
	"testing"

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

	fnRuntime := NewChannel[*TestData](loop)

	startErr := fnRuntime.Start()

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

	fnRuntime.Stop()

	fnRuntime.(*Runner).wait.Wait()

	assert.NoError(startErr)
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

	fnRuntime := NewChannel[*TestData](loop)

	startErr := fnRuntime.Start()
	assert.Error(startErr)
	fnRuntime.(*Runner).wait.Wait()

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

	fnRuntime := NewChannel[*TestData](loop)

	startErr := fnRuntime.Start()
	assert.NoError(startErr)
	fnRuntime.(*Runner).wait.Wait()

	assert.Equal(0, value)
}

func TestChannelNilChannel(t *testing.T) {
	assert := assert.New(t)

	value := 0

	fnRuntime := NewChannel[*TestData](nil)

	startErr := fnRuntime.Start()
	assert.Error(startErr)
	fnRuntime.(*Runner).wait.Wait()

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

	fnRuntime := NewChannel[*TestData](loop)

	startErr := fnRuntime.Start()

	loop.channel <- &TestData{}
	loop.channel <- &TestData{}
	loop.channel <- &TestData{}

	fnRuntime.(*Runner).wait.Wait()

	assert.NoError(startErr)
	assert.Equal(3, value)
}

type TestChannel struct {
	brokenInit    bool
	brokenChannel bool
	channel       chan *TestData
	loop          func(*TestData) error
}

func (l *TestChannel) Start() error {
	l.channel = make(chan *TestData)
	if l.brokenInit {
		return errors.New("broken init")
	}
	return nil
}

func (l *TestChannel) Stop() {
	close(l.channel)
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
