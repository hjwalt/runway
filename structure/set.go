package structure

import "sync"

type Set[T comparable] interface {
	Add(...T)
	Contain(...T) bool
	Remove(...T)
	Clear()
	Get() []T
}

func NewSet[T comparable]() Set[T] {
	return &mapSet[T]{
		internal: map[T]bool{},
	}
}

func NewSetFrom[T comparable](vals []T) Set[T] {
	setMap := &mapSet[T]{
		internal: map[T]bool{},
	}
	for _, v := range vals {
		setMap.Add(v)
	}
	return setMap
}

type mapSet[T comparable] struct {
	internal map[T]bool
	sync     sync.Mutex
}

func (s *mapSet[T]) Add(ts ...T) {
	s.sync.Lock()
	defer s.sync.Unlock()
	for _, t := range ts {
		s.internal[t] = true
	}
}

func (s *mapSet[T]) Contain(ts ...T) bool {
	s.sync.Lock()
	defer s.sync.Unlock()
	for _, t := range ts {
		if _, tPresent := s.internal[t]; !tPresent {
			return false
		}
	}
	return true
}

func (s *mapSet[T]) Remove(ts ...T) {
	s.sync.Lock()
	defer s.sync.Unlock()
	for _, t := range ts {
		delete(s.internal, t)
	}
}

func (s *mapSet[T]) Clear() {
	s.sync.Lock()
	defer s.sync.Unlock()
	s.internal = map[T]bool{}
}

func (s *mapSet[T]) Get() []T {
	s.sync.Lock()
	defer s.sync.Unlock()
	result := make([]T, len(s.internal))
	i := 0
	for t := range s.internal {
		result[i] = t
		i += 1
	}
	return result
}
