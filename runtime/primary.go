package runtime

import (
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/hjwalt/runway/logger"
)

// constructor
func NewPrimary(configurations ...Configuration[*Primary]) Runtime {
	c := &Primary{}
	c = PrimaryDefault(c)
	for _, configuration := range configurations {
		c = configuration(c)
	}
	return c
}

// default
func PrimaryDefault(c *Primary) *Primary {
	controller, err := NewPrimaryController()
	c.SetController(controller)
	c.err = err
	return c
}

// configuration
func PrimaryWithRuntime(runtime Runtime) Configuration[*Primary] {
	return func(c *Primary) *Primary {
		if c.runtimes == nil {
			c.runtimes = make([]Runtime, 0)
		}
		runtime.SetController(c.controller)
		c.runtimes = append(c.runtimes, runtime)
		return c
	}
}

// implementation
type Primary struct {
	controller Controller
	err        chan error
	runtimes   []Runtime
}

func (r *Primary) Start() error {
	var startError error

	for i := 0; i < len(r.runtimes); i++ {
		if err := r.runtimes[i].Start(); err != nil {
			r.StopFrom(i - 1)
			startError = err
			break
		}
	}

	if startError != nil {
		r.controller.Wait()
		return errors.Join(ErrPrimaryInitialiseError, startError)
	}

	go r.Interrupt()
	go r.Error()

	r.controller.Wait()

	return nil
}

func (r *Primary) Stop() {
	r.StopFrom(len(r.runtimes) - 1)
}

func (r *Primary) SetController(controller Controller) {
	r.controller = controller
}

func (r *Primary) StopFrom(from int) {
	for i := from; i > -1; i-- {
		r.runtimes[i].Stop()
	}
}

func (r *Primary) Interrupt() {
	interruptSignal := make(chan os.Signal, 10)
	signal.Notify(interruptSignal, os.Interrupt, syscall.SIGTERM)
	<-interruptSignal
	r.Stop()
}

func (r *Primary) Error() {
	err := <-r.err
	logger.ErrorErr("runtime error received, exitting", err)
	r.Stop()
	if !errors.Is(err, ErrPrimaryTesting) {
		os.Exit(1) // will fail the unit test, so have to flag out
	}
}

// Errors
var (
	ErrPrimaryInitialiseError = errors.New("primary runtime initialise function failed")
	ErrPrimaryTesting         = errors.New("primary runtime testing error")
)
