package inverse

import (
	"context"

	"github.com/hjwalt/runway/reflect"
)

type Injector[T any] func(context.Context, Container) (T, error)

func GenericAdd[T any](c Container, qualifier string, injector Injector[T]) {
	c.Add(qualifier, func(ctx context.Context, ic Container) (any, error) { return injector(ctx, ic) })
}

func GenericAddVal[T any](c Container, qualifier string, instance T) {
	c.AddVal(qualifier, instance)
}

func GenericGet[T any](c Container, ctx context.Context, qualifier string) (T, error) {
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

// Deprecated: use GenericGet instead to keep naming scheme consistent
func GenericGetLast[T any](c Container, ctx context.Context, qualifier string) (T, error) {
	return GenericGet[T](c, ctx, qualifier)
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
