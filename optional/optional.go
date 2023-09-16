package optional

import "github.com/hjwalt/runway/reflect"

func Of[T any](v T) Optional[T] {
	return optional[T]{internal: &v}
}

func OfPointer[T any](v *T) Optional[T] {
	return optional[T]{internal: v}
}

type Optional[T any] interface {
	IsPresent() bool
	Get() T
	GetOrDefault(defaultValue T) T
}

type optional[T any] struct {
	internal *T
}

func (o optional[T]) IsPresent() bool {
	return o.internal != nil
}

func (o optional[T]) Get() T {
	if o.internal == nil {
		return reflect.Construct[T]()
	} else {
		return *o.internal
	}
}

func (o optional[T]) GetOrDefault(defaultValue T) T {
	if o.internal == nil {
		return defaultValue
	} else {
		return *o.internal
	}
}
