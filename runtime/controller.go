package runtime

import "sync"

func NewController() Controller {
	return &RuntimeController{}
}

type RuntimeController struct {
	wait sync.WaitGroup
}

func (c *RuntimeController) Started() {
	c.wait.Add(1)
}

func (c *RuntimeController) Stopped() {
	c.wait.Add(-1)
}

func (c *RuntimeController) Wait() {
	c.wait.Wait()
}
