package inverse

import (
	"context"
)

func GetAll[T any](ctx context.Context, qualifier string) ([]T, error) {
	cachedInstances, cachedExist := qualifierInjectedAll[qualifier]
	if cachedExist {
		results := make([]T, len(cachedInstances))
		for i, cached := range cachedInstances {
			cachedCasted, castedOk := cached.(T)
			if !castedOk {
				return results, ErrorCastingFailure(qualifier)
			}
			results[i] = cachedCasted
		}
		return results, nil
	}

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

	anyArray := make([]any, len(resultList))
	for i, resAny := range resultList {
		anyArray[i] = resAny
	}
	qualifierInjectedAll[qualifier] = anyArray

	return resultList, nil
}
