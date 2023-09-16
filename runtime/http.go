package runtime

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/hjwalt/runway/logger"
)

// constructor
var NewHttp = ConstructorFor[*HttpRunnable, Runtime](
	func() *HttpRunnable {
		return &HttpRunnable{
			server: &http.Server{
				Addr:         ":8080",
				ReadTimeout:  5 * time.Second,
				WriteTimeout: 5 * time.Second,
			},
		}
	},
	func(hr *HttpRunnable) Runtime {
		return NewRunner(hr)
	},
)

// configuration
func HttpWithHandler(handler http.Handler) Configuration[*HttpRunnable] {
	return func(c *HttpRunnable) *HttpRunnable {
		c.server.Handler = handler
		return c
	}
}

func HttpWithPort(port int) Configuration[*HttpRunnable] {
	return func(c *HttpRunnable) *HttpRunnable {
		c.server.Addr = ":" + fmt.Sprintf("%d", port)
		return c
	}
}

func HttpWithReadTimeout(timeout time.Duration) Configuration[*HttpRunnable] {
	return func(c *HttpRunnable) *HttpRunnable {
		c.server.ReadTimeout = timeout
		return c
	}
}

func HttpWithReadHeaderTimeout(timeout time.Duration) Configuration[*HttpRunnable] {
	return func(c *HttpRunnable) *HttpRunnable {
		c.server.ReadHeaderTimeout = timeout
		return c
	}
}

func HttpWithWriteTimeout(timeout time.Duration) Configuration[*HttpRunnable] {
	return func(c *HttpRunnable) *HttpRunnable {
		c.server.WriteTimeout = timeout
		return c
	}
}

// implementation
type HttpRunnable struct {
	server *http.Server
}

func (c *HttpRunnable) Start() error {
	if c.server == nil {
		return ErrHttpMissingServer
	}
	if c.server.Handler == nil {
		return ErrHttpMissingHandler
	}
	return nil
}

func (c *HttpRunnable) Stop() {
	if c.server == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := c.server.Shutdown(ctx)
	logger.ErrorIfErr("Server forced to shutdown", err)
}

func (c *HttpRunnable) Run() error {
	err := c.server.ListenAndServe()
	if err != nil && err == http.ErrServerClosed {
		err = nil
	}
	return err
}

// Errors
var (
	ErrHttpMissingServer  = errors.New("http runtime no server provided")
	ErrHttpMissingHandler = errors.New("http runtime no handler provided")
)
