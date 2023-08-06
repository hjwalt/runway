package inverse

import (
	"context"
	"sync"
)

type Injector func(context.Context) (any, error)

var qualifierMutex = sync.Mutex{}

var qualifierInjectorMap = map[string][]Injector{}

var qualifierInjectedLast = map[string]any{}

func Reset() {
	qualifierMutex.Lock()
	defer qualifierMutex.Unlock()

	qualifierInjectorMap = map[string][]Injector{}
	Release()
}

func Release() {
	qualifierInjectedLast = map[string]any{}
}

func Register(qualifier string, injector Injector) {
	qualifierMutex.Lock()
	defer qualifierMutex.Unlock()

	qualifierList, qualifierExist := qualifierInjectorMap[qualifier]
	if !qualifierExist {
		qualifierList = make([]Injector, 0)
	}

	qualifierList = append(qualifierList, injector)
	qualifierInjectorMap[qualifier] = qualifierList
}
