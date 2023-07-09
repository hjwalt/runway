package format

// bytes format
type Format[T any] interface {
	Default() T
	Marshal(T) ([]byte, error)
	Unmarshal([]byte) (T, error)
}
