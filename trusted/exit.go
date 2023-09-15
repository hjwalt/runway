package trusted

import (
	"errors"
	"os"
)

// Errors
var (
	ErrPrimaryTesting = errors.New("primary runtime testing error")
)

func Exit(err error) {
	if !errors.Is(err, ErrPrimaryTesting) {
		os.Exit(1) // will fail the unit test, so have to flag out
	}
}
