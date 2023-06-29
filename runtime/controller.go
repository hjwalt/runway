package runtime

import "sync"

func NewPrimaryController() (Controller, chan error) {
	err := make(chan error, 1)
	return &RuntimeController{err: err}, err
}

type RuntimeController struct {
	wait sync.WaitGroup
	err  chan<- error
}

func (c *RuntimeController) Started() {
	c.wait.Add(1)
}

func (c *RuntimeController) Stopped() {
	c.wait.Add(-1)
}

func (c *RuntimeController) Error(err error) {
	c.err <- err
}

func (c *RuntimeController) Wait() {
	c.wait.Wait()
}
