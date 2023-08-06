package inverse

import (
	"context"
)

func GetAll[T any](ctx context.Context, qualifier string) ([]T, error) {
	// Check for context and loop
	if ctx == nil {
		return make([]T, 0), ErrorNilContext(qualifier)
	}
	if ctx.Value(Qualifier(qualifier)) != nil {
		return make([]T, 0), ErrorResolveLoop(qualifier)
	}

	// Resolve
	qualifierMutex.Lock()
	qualifierList, qualifierExist := qualifierInjectorMap[qualifier]
	qualifierMutex.Unlock()

	resultList := make([]T, 0)
	if !qualifierExist {
		return resultList, ErrorNotInjected(qualifier)
	}

	outerCastOk := true

	for _, injector := range qualifierList {
		injectedVal, err := injector(context.WithValue(ctx, Qualifier(qualifier), "already-resolved"))
		if err != nil {
			return resultList, err
		}

		castedVal, castOk := injectedVal.(T)

		outerCastOk = outerCastOk && castOk
		if castOk {
			resultList = append(resultList, castedVal)
		}
	}

	if !outerCastOk {
		return resultList, ErrorCastingFailure(qualifier)
	}

	return resultList, nil
}
