package structure

import "sync"

type MultiMap[K comparable, V any] interface {
	Add(K, ...V)
	Get(K) []V
	Contain(K) bool
	Remove(K)
	Clear()
}

func NewMultiMap[K comparable, V any]() MultiMap[K, V] {
	return &arrayMultiMap[K, V]{
		internal: map[K][]V{},
	}
}

type arrayMultiMap[K comparable, V any] struct {
	internal map[K][]V
	sync     sync.Mutex
}

func (mm *arrayMultiMap[K, V]) Add(k K, vs ...V) {
	mm.sync.Lock()
	defer mm.sync.Unlock()
	if currVals, kExist := mm.internal[k]; kExist {
		mm.internal[k] = append(currVals, vs...)
	} else {
		mm.internal[k] = append([]V{}, vs...)
	}
}

func (mm *arrayMultiMap[K, V]) Get(k K) []V {
	mm.sync.Lock()
	defer mm.sync.Unlock()
	if currVals, kExist := mm.internal[k]; kExist {
		return currVals
	} else {
		return []V{}
	}
}

func (mm *arrayMultiMap[K, V]) Contain(k K) bool {
	mm.sync.Lock()
	defer mm.sync.Unlock()
	if currVals, kExist := mm.internal[k]; kExist {
		return len(currVals) > 0
	} else {
		return false
	}
}

func (mm *arrayMultiMap[K, V]) Remove(k K) {
	mm.sync.Lock()
	defer mm.sync.Unlock()
	delete(mm.internal, k)
}

func (mm *arrayMultiMap[K, V]) Clear() {
	mm.sync.Lock()
	defer mm.sync.Unlock()
	mm.internal = map[K][]V{}
}
