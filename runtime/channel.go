package runtime

import (
	"errors"

	"github.com/hjwalt/runway/logger"
)

type Channel[T any] interface {
	Start() error
	Stop()
	Channel() (<-chan T, error)
	Loop(T) error
}

// constructor
func NewChannel[T any](channel Channel[T]) Runtime {
	c := &ChannelRunnable[T]{
		channel: channel,
	}

	return NewRunner(c)
}

// implementation
type ChannelRunnable[T any] struct {
	channel Channel[T]
}

func (r *ChannelRunnable[T]) Start() error {
	if r.channel == nil {
		return ErrChannelRuntimeNoChannel
	}

	initerr := r.channel.Start()
	if initerr != nil {
		return errors.Join(ErrChannelRuntimeInitialise, initerr)
	}

	return nil
}

func (r *ChannelRunnable[T]) Stop() {
	if r.channel != nil {
		r.channel.Stop()
	}
}

func (r *ChannelRunnable[T]) Run() error {
	channel, channelErr := r.channel.Channel()
	if channelErr != nil {
		return channelErr
	}

	for v := range channel {
		err := r.channel.Loop(v)
		if err != nil {
			logger.ErrorErr("functional runtime Channel error", err)
			return err
		}
	}

	return nil
}

// Errors
var (
	ErrChannelRuntimeNoChannel  = errors.New("functional runtime no Channel function provided")
	ErrChannelRuntimeInitialise = errors.New("functional runtime initialise function failed")
)
