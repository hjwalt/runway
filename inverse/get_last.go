package inverse

import (
	"context"

	"github.com/hjwalt/runway/logger"
	"github.com/hjwalt/runway/reflect"
)

func GetLast[T any](ctx context.Context, qualifier string) (T, error) {
	cached, cachedExist := qualifierInjectedLast[qualifier]
	if cachedExist {
		cachedCasted, castedOk := cached.(T)
		if !castedOk {
			return reflect.Construct[T](), ErrorCastingFailure(qualifier)
		}
		return cachedCasted, nil
	}

	// Check for context and loop
	if ctx == nil {
		return reflect.Construct[T](), ErrorNilContext(qualifier)
	}
	if ctx.Value(Qualifier(qualifier)) != nil {
		logger.Info("qualifier loop")
		return reflect.Construct[T](), ErrorResolveLoop(qualifier)
	}

	// Resolve
	qualifierMutex.Lock()
	qualifierList, qualifierExist := qualifierInjectorMap[qualifier]
	qualifierMutex.Unlock()

	if !qualifierExist {
		return reflect.Construct[T](), ErrorNotInjected(qualifier)
	}

	lastInjector := qualifierList[len(qualifierList)-1]
	injectedVal, err := lastInjector(context.WithValue(ctx, Qualifier(qualifier), "already-resolved"))
	if err != nil {
		return reflect.Construct[T](), err
	}

	castedVal, castOk := injectedVal.(T)
	if !castOk {
		return reflect.Construct[T](), ErrorCastingFailure(qualifier)
	}

	qualifierInjectedLast[qualifier] = castedVal

	return castedVal, nil
}
