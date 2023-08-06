package inverse

import "context"

func RegisterInstance[T any](qualifier string, instance T) {
	Register[T](qualifier, func(ctx context.Context) (T, error) { return instance, nil })
}

func RegisterInstances[T any](qualifier string, instances []T) {
	for _, instance := range instances {
		RegisterInstance[T](qualifier, instance)
	}
}
