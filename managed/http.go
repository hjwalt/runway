package managed

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/hjwalt/runway/inverse"
)

const (
	ConfHttpPort                    = "ConfHttpPort"
	ConfHttpReadTimeoutMillisecond  = "ConfReadTimeoutMillisecond"
	ConfHttpWriteTimeoutMillisecond = "ConfWriteTimeoutMillisecond"
	ConfHttpTlsKeyPath              = "ConfHttpTlsKeyPath"
	ConfHttpTlsCertPath             = "ConfHttpTlsCertPath"
)

func NewHttp() []Service {
	httpService := &httpRunnable{}
	return []Service{httpService, NewRunner(httpService)}
}

// implementation
type httpRunnable struct {
	server      *http.Server
	tlsKeyPath  string
	tlsCertPath string
}

func (c *httpRunnable) Name() string {
	return "http"
}

func (r *httpRunnable) Register(ctx context.Context, ic inverse.Container) error {
	return nil
}

func (r *httpRunnable) Resolve(ctx context.Context, ic inverse.Container) error {
	config, configErr := ResolveConfig(ctx, ic, r.Name())
	if configErr != nil {
		return configErr
	}

	port := config.GetString(ConfHttpPort, "8080")

	r.server = &http.Server{
		Addr:         ":" + port,
		ReadTimeout:  time.Duration(config.GetInt64(ConfHttpReadTimeoutMillisecond, 5000)) * time.Millisecond,
		WriteTimeout: time.Duration(config.GetInt64(ConfHttpWriteTimeoutMillisecond, 5000)) * time.Millisecond,
	}

	handler, handlerErr := inverse.GenericGet[http.Handler](ic, ctx, QualifierHttpHandler)
	if handlerErr != nil {
		return handlerErr
	}
	r.server.Handler = handler

	r.tlsCertPath = config.GetString(ConfHttpTlsCertPath, "")
	r.tlsKeyPath = config.GetString(ConfHttpTlsKeyPath, "")

	return nil
}

func (r *httpRunnable) Clean() error {
	return nil
}

func (r *httpRunnable) Start() error {

	if r.server == nil {
		return ErrHttpMissingServer
	}
	if r.server.Handler == nil {
		return ErrHttpMissingHandler
	}

	return nil
}

func (r *httpRunnable) Stop() error {
	if r.server == nil {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return r.server.Shutdown(ctx)
}

func (c *httpRunnable) Run() error {
	var err error

	if len(c.tlsCertPath) > 0 && len(c.tlsKeyPath) > 0 {
		err = c.server.ListenAndServeTLS(c.tlsCertPath, c.tlsKeyPath)
	} else {
		err = c.server.ListenAndServe()
	}
	if err != nil && err == http.ErrServerClosed {
		err = nil
	}
	return err
}

// Errors
var (
	ErrHttpMissingServer         = errors.New("http runtime no server provided")
	ErrHttpMissingHandler        = errors.New("http runtime no handler provided")
	ErrHttpFailedToInitialiseTls = errors.New("http runtime failed to start tls")
)
