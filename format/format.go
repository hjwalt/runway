package format

// bytes format
type Format[T any] interface {
	Default() T
	Marshal(T) ([]byte, error)
	Unmarshal([]byte) (T, error)
}

type Mask interface {
	Mask([]byte) ([]byte, error)
	Unmask([]byte) ([]byte, error)
}
