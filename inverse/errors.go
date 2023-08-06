package inverse

import (
	"errors"
	"fmt"
)

var ErrNotInjected = errors.New("qualifier not injected")

func ErrorNilContext(qualifier string) error {
	return fmt.Errorf("context is nil while resolving qualifier %s", qualifier)
}

func ErrorResolveLoop(qualifier string) error {
	return fmt.Errorf("loop detected for qualifier %s", qualifier)
}

func ErrorNotInjected(qualifier string) error {
	return errors.Join(ErrNotInjected, fmt.Errorf("qualifier %s not injected", qualifier))
}

func ErrorCastingFailure(qualifier string) error {
	return fmt.Errorf("qualifier %s casting failed", qualifier)
}
