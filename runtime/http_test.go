package runtime_test

import (
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/hjwalt/runway/runtime"
	"github.com/stretchr/testify/assert"
)

func TestHttp(t *testing.T) {
	assert := assert.New(t)

	controller := NewController()

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "This is my website!\n")
	})

	httpRuntime := runtime.NewHttp(
		runtime.HttpWithPort(30080),
		runtime.HttpWithHandler(mux),
		runtime.HttpWithReadTimeout(5*time.Second),
		runtime.HttpWithReadHeaderTimeout(5*time.Second),
		runtime.HttpWithWriteTimeout(5*time.Second),
	)
	httpRuntime.SetController(controller)

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
	controller.Wait()
}
