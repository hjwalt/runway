package runtime

import (
	"errors"
	"os"
	"os/signal"
	"syscall"
)

// constructor
func NewPrimary(configurations ...Configuration[*Primary]) Runtime {
	consumer := &Primary{}
	for _, configuration := range configurations {
		consumer = configuration(consumer)
	}
	return consumer
}

// configuration
func PrimaryWithController(controller Controller) Configuration[*Primary] {
	return func(c *Primary) *Primary {
		c.controller = controller
		return c
	}
}

func PrimaryWithRuntime(runtime Runtime) Configuration[*Primary] {
	return func(c *Primary) *Primary {
		if c.runtimes == nil {
			c.runtimes = make([]Runtime, 0)
		}
		c.runtimes = append(c.runtimes, runtime)
		return c
	}
}

// implementation
type Primary struct {
	controller Controller
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
	// defer r.Panic()

	r.controller.Wait()

	return nil
}

func (r *Primary) StopFrom(from int) {
	for i := from; i > -1; i-- {
		r.runtimes[i].Stop()
	}
}

func (r *Primary) Stop() {
	r.StopFrom(len(r.runtimes) - 1)
}

func (runtime *Primary) Interrupt() {
	interruptSignal := make(chan os.Signal, 10)
	signal.Notify(interruptSignal, os.Interrupt, syscall.SIGTERM)
	<-interruptSignal
	runtime.Stop()
	// os.Exit(0) -- don't need os.Exit if everything is cleaned up properly
}

// Deal with panic some other time
// func (runtime *Primary) Panic() {
// 	logger.Infof("panicking")
// 	if x := recover(); x != nil {
// 		logger.Error("runtime panic", zap.Error(x.(error)))
// 		runtime.Stop()
// 		os.Exit(1)
// 	}
// }

// Errors
var (
	ErrPrimaryInitialiseError = errors.New("primary runtime initialise function failed")
)
