package runtime

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunnerMissingRunnable(t *testing.T) {
	assert := assert.New(t)

	runner := NewRunner(nil)
	err := runner.Start()
	assert.Equal(ErrRunnerRuntimeNoRunnable, err)
}
