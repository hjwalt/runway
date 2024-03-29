package configuration

import (
	"errors"
	"os"

	"github.com/hjwalt/runway/format"
)

func Read[T any](file string, f format.Format[T]) (T, error) {
	bytes, err := os.ReadFile(file)
	if err != nil {
		return f.Default(), errors.Join(ErrReadFail, err)
	}
	return f.Unmarshal(bytes)
}

var ErrReadFail = errors.New("cannot read file")
