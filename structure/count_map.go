package structure

import "sync"

type CountMap[K comparable] interface {
	Set(K, int64)
	Add(K, int64)
	Subtract(K, int64)
	Get(K) int64
	GetAll() map[K]int64
	Contain(K) bool
	Remove(K)
	Clear()
}

func NewCountMap[K comparable]() CountMap[K] {
	return &basicCountMap[K]{
		internal: map[K]int64{},
	}
}

type basicCountMap[K comparable] struct {
	internal map[K]int64
	sync     sync.Mutex
}

func (mm *basicCountMap[K]) Set(k K, vs int64) {
	mm.sync.Lock()
	defer mm.sync.Unlock()
	mm.internal[k] = vs
}

func (mm *basicCountMap[K]) Add(k K, vs int64) {
	mm.sync.Lock()
	defer mm.sync.Unlock()
	if currVals, kExist := mm.internal[k]; kExist {
		mm.internal[k] = currVals + vs
	} else {
		mm.internal[k] = vs
	}
}

func (mm *basicCountMap[K]) Subtract(k K, vs int64) {
	mm.sync.Lock()
	defer mm.sync.Unlock()
	if currVals, kExist := mm.internal[k]; kExist {
		mm.internal[k] = currVals - vs
	} else {
		mm.internal[k] = -vs
	}
}

func (mm *basicCountMap[K]) Get(k K) int64 {
	mm.sync.Lock()
	defer mm.sync.Unlock()
	if currVals, kExist := mm.internal[k]; kExist {
		return currVals
	} else {
		return 0
	}
}

func (mm *basicCountMap[K]) GetAll() map[K]int64 {
	mm.sync.Lock()
	defer mm.sync.Unlock()
	duplicatedMap := map[K]int64{}
	for k, v := range mm.internal {
		duplicatedMap[k] = v
	}
	return duplicatedMap
}

func (mm *basicCountMap[K]) Contain(k K) bool {
	mm.sync.Lock()
	defer mm.sync.Unlock()
	_, kExist := mm.internal[k]
	return kExist
}

func (mm *basicCountMap[K]) Remove(k K) {
	mm.sync.Lock()
	defer mm.sync.Unlock()
	delete(mm.internal, k)
}

func (mm *basicCountMap[K]) Clear() {
	mm.sync.Lock()
	defer mm.sync.Unlock()
	mm.internal = map[K]int64{}
}
