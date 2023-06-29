package runtime

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hjwalt/runway/logger"
	"go.uber.org/zap"
)

// constructor
func NewHttp(configurations ...Configuration[*Http]) Runtime {
	c := &Http{}
	c = HttpDefault(c)
	for _, configuration := range configurations {
		c = configuration(c)
	}
	return c
}

// default
func HttpDefault(c *Http) *Http {
	c.server = &http.Server{
		Addr:         ":8080",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	return c
}

// configuration
func HttpWithHandler(handler http.Handler) Configuration[*Http] {
	return func(c *Http) *Http {
		c.server.Handler = handler
		return c
	}
}

func HttpWithPort(port int) Configuration[*Http] {
	return func(c *Http) *Http {
		c.server.Addr = ":" + fmt.Sprintf("%d", port)
		return c
	}
}

func HttpWithReadTimeout(timeout time.Duration) Configuration[*Http] {
	return func(c *Http) *Http {
		c.server.ReadTimeout = timeout
		return c
	}
}

func HttpWithReadHeaderTimeout(timeout time.Duration) Configuration[*Http] {
	return func(c *Http) *Http {
		c.server.ReadHeaderTimeout = timeout
		return c
	}
}

func HttpWithWriteTimeout(timeout time.Duration) Configuration[*Http] {
	return func(c *Http) *Http {
		c.server.WriteTimeout = timeout
		return c
	}
}

// implementation
type Http struct {
	controller Controller
	server     *http.Server
}

func (c *Http) Start() error {
	go c.Run()
	c.controller.Started()
	return nil
}

func (c *Http) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := c.server.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown: ", zap.Error(err))
	}
}

func (c *Http) SetController(controller Controller) {
	c.controller = controller
}

func (c *Http) Run() {
	defer c.controller.Stopped()
	if err := c.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		c.controller.Error(err)
	}
}
