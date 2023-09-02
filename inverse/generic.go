package inverse

import (
	"context"

	"github.com/hjwalt/runway/reflect"
)

func GenericAdd[T any](c Container, qualifier string, injector Injector[T]) {
	c.Add(qualifier, func(ctx context.Context) (any, error) { return injector(ctx) })
}

func GenericGetLast[T any](c Container, ctx context.Context, qualifier string) (T, error) {
	injectedVal, err := c.Get(ctx, qualifier)
	if err != nil {
		return reflect.Construct[T](), err
	}
	if castedVal, castOk := injectedVal.(T); castOk {
		return castedVal, nil
	} else {
		c.Invalid(ErrInverseCastError)
		return reflect.Construct[T](), getError(qualifier, ErrInverseCastError)
	}
}

func GenericGetAll[T any](c Container, ctx context.Context, qualifier string) ([]T, error) {
	vals, err := c.GetAll(ctx, qualifier)
	if err != nil {
		return []T{}, err
	}
	resultList := make([]T, len(vals))
	for resi, injectedVal := range vals {
		if castedVal, castOk := injectedVal.(T); castOk {
			resultList[resi] = castedVal
		} else {
			c.Invalid(ErrInverseCastError)
			return []T{}, getError(qualifier, ErrInverseCastError)
		}
	}
	return resultList, nil
}
