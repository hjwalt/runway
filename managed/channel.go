package managed

import "errors"

func NewChannel[T any](channel Channel[T]) Service {
	c := &channelRunnable[T]{
		channel: channel,
	}
	return NewRunner(c)
}

type channelRunnable[T any] struct {
	channel Channel[T]
}

func (r *channelRunnable[T]) Run() error {
	if r.channel == nil {
		return ErrChannelRuntimeNoChannel
	}

	channel, channelErr := r.channel.Channel()
	if channelErr != nil {
		return channelErr
	}

	for v := range channel {
		err := r.channel.Loop(v)
		if err != nil {
			return err
		}
	}

	return nil
}

// Errors
var (
	ErrChannelRuntimeNoChannel = errors.New("functional runtime no Channel function provided")
)
