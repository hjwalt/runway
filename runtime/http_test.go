package runtime

import (
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHttpShouldStopNormally(t *testing.T) {
	assert := assert.New(t)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "This is my website!\n")
	})

	httpRuntime := NewHttp(
		HttpWithPort(30080),
		HttpWithHandler(mux),
		HttpWithReadTimeout(5*time.Second),
		HttpWithReadHeaderTimeout(5*time.Second),
		HttpWithWriteTimeout(5*time.Second),
	)

	startErr := httpRuntime.Start()
	assert.NoError(startErr)

	var resp *http.Response
	var geterr error

	for i := 0; i < 10; i++ {
		resp, geterr = http.Get("http://localhost:30080")
		if geterr == nil {
			break
		}
		time.Sleep(1 * time.Millisecond)
	}
	assert.NoError(geterr)

	defer resp.Body.Close()
	body, readerr := io.ReadAll(resp.Body)
	assert.NoError(readerr)

	assert.Equal("This is my website!\n", string(body))

	httpRuntime.Stop()
	httpRuntime.Stop() // testing multiple stop will not break the system
}

func TestHttpMissingServer(t *testing.T) {
	assert := assert.New(t)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "This is my website!\n")
	})

	httpRuntime := HttpRunnable{}

	startErr := httpRuntime.Start()
	assert.ErrorIs(startErr, ErrHttpMissingServer)
	httpRuntime.Stop()
}

func TestHttpMissingHandler(t *testing.T) {
	assert := assert.New(t)

	httpRuntime := NewHttp(
		HttpWithPort(30080),
		HttpWithReadTimeout(5*time.Second),
		HttpWithReadHeaderTimeout(5*time.Second),
		HttpWithWriteTimeout(5*time.Second),
	)

	startErr := httpRuntime.Start()
	assert.ErrorIs(startErr, ErrHttpMissingHandler)
	httpRuntime.Stop()
}

func TestHttpsShouldSucceed(t *testing.T) {
	assert := assert.New(t)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "This is my website!\n")
	})

	httpRuntime := NewHttp(
		HttpWithPort(30080),
		HttpWithHandler(mux),
		HttpWithReadTimeout(5*time.Second),
		HttpWithReadHeaderTimeout(5*time.Second),
		HttpWithWriteTimeout(5*time.Second),
		HttpWithTls("./tls/server.crt", "./tls/server.key"),
	)

	startErr := httpRuntime.Start()
	assert.NoError(startErr)

	var resp *http.Response
	var geterr error

	for i := 0; i < 10; i++ {
		resp, geterr = http.Get("http://localhost:30080")
		if geterr == nil {
			break
		}
		time.Sleep(1 * time.Millisecond)
	}
	assert.NoError(geterr)

	defer resp.Body.Close()
	body, readerr := io.ReadAll(resp.Body)
	assert.NoError(readerr)

	assert.Equal("This is my website!\n", string(body))

	httpRuntime.Stop()
}

func TestHttpsShouldFailToStart(t *testing.T) {
	assert := assert.New(t)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "This is my website!\n")
	})

	httpRuntime := NewHttp(
		HttpWithPort(30080),
		HttpWithHandler(mux),
		HttpWithReadTimeout(5*time.Second),
		HttpWithReadHeaderTimeout(5*time.Second),
		HttpWithWriteTimeout(5*time.Second),
		HttpWithTls("./tls/server.crt", "./tls/wrong.key"),
	)

	startErr := httpRuntime.Start()
	assert.ErrorIs(startErr, ErrHttpFailedToInitialiseTls)
	httpRuntime.Stop()
}
