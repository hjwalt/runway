package runtime

import (
	"errors"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/hjwalt/runway/logger"
	"github.com/hjwalt/runway/trusted"
)

var primary atomic.Pointer[Primary]

type Runtime interface {
	Start() error
	Stop()
}

func Start(runtimes []Runtime, sleepTime time.Duration) error {
	c := &Primary{
		runtimes:         runtimes,
		interruptChannel: make(chan os.Signal, 10),
		errorChannel:     make(chan error, 10),
		sleepTime:        sleepTime,
	}
	primary.Store(c)
	return primary.Load().Start()
}

func Stop() {
	if primary.Load() != nil {
		primary.Load().Stop()
	}
}

func Wait() {
	if primary.Load() != nil {
		primary.Load().Wait()
	}
}

func Error(err error) {
	if primary.Load() != nil {
		primary.Load().errorChannel <- err
	}
}

// implementation
type Primary struct {
	// wait             sync.WaitGroup
	started          atomic.Bool
	runtimes         []Runtime
	interruptChannel chan os.Signal
	errorChannel     chan error
	sleepTime        time.Duration
}

func (r *Primary) Start() error {
	logger.Info("starting up")

	var startError error

	for i := 0; i < len(r.runtimes); i++ {
		if err := r.runtimes[i].Start(); err != nil {
			r.StopFrom(i - 1)
			startError = err
			break
		}
	}

	if startError != nil {
		return errors.Join(ErrPrimaryInitialiseError, startError)
	}

	go r.Interrupt()
	go r.Error()

	r.started.Store(true)
	return nil
}

func (r *Primary) Stop() {
	defer r.started.Store(false)
	logger.Info("shutting down")
	r.StopFrom(len(r.runtimes) - 1)
}

func (r *Primary) Wait() {
	for r.started.Load() {
		time.Sleep(r.sleepTime)
	}
}

func (r *Primary) StopFrom(from int) {
	for i := from; i > -1; i-- {
		r.runtimes[i].Stop()
	}
}

func (r *Primary) Interrupt() {
	signal.Notify(r.interruptChannel, os.Interrupt, syscall.SIGTERM)
	<-r.interruptChannel
	r.Stop()
}

func (r *Primary) Error() {
	err := <-r.errorChannel
	logger.ErrorErr("runtime error received, exitting", err)
	r.Stop()
	trusted.Exit(err)
}

// Errors
var (
	ErrPrimaryInitialiseError = errors.New("primary runtime initialise function failed")
)
