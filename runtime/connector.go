package runtime

import (
	"errors"

	"github.com/hjwalt/runway/reflect"
)

// constructor
func NewConnector[C any](configurations ...Configuration[*Connector[C]]) Runtime {
	c := &Connector[C]{}
	c = ConnectorDefault[C](c)
	for _, configuration := range configurations {
		c = configuration(c)
	}
	return c
}

// default
func ConnectorDefault[C any](c *Connector[C]) *Connector[C] {
	c.data = reflect.Construct[C]()
	return c
}

// configuration
func ConnectorWithInitialise[C any](initialise func() (C, error)) Configuration[*Connector[C]] {
	return func(c *Connector[C]) *Connector[C] {
		c.initialise = initialise
		return c
	}
}

func ConnectorWithCleanup[C any](cleanup func(C)) Configuration[*Connector[C]] {
	return func(c *Connector[C]) *Connector[C] {
		c.cleanup = cleanup
		return c
	}
}

// implementation
type Connector[C any] struct {
	controller Controller
	initialise func() (C, error)
	cleanup    func(C)
	data       C
}

func (r *Connector[C]) Start() error {
	if r.initialise != nil {
		data, initerr := r.initialise()
		if initerr != nil {
			return errors.Join(ErrConnectorRuntimeInitialise, initerr)
		}
		r.data = data
	}

	r.controller.Started()
	return nil
}

func (r *Connector[C]) Stop() {

	if r.cleanup != nil {
		r.cleanup(r.data)
	}

	r.controller.Stopped()
}

func (r *Connector[C]) SetController(controller Controller) {
	r.controller = controller
}

// Errors
var (
	ErrConnectorRuntimeInitialise = errors.New("connector runtime initialise function failed")
)
