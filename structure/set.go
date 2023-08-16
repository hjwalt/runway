package structure

import "sync"

type Set[T comparable] interface {
	Add(T)
	AddAll(...T)
}

func NewSet[T comparable]() Set[T] {
	return &MapSet[T]{}
}

type MapSet[T comparable] struct {
	internal map[T]bool
	sync     sync.Mutex
}

func (s *MapSet[T]) Add(t T) {
	s.sync.Lock()
	defer s.sync.Unlock()
	s.internal[t] = true
}

func (s *MapSet[T]) AddAll(ts ...T) {
	s.sync.Lock()
	defer s.sync.Unlock()
	for _, t := range ts {
		s.internal[t] = true
	}
}

func (s *MapSet[T]) Contain(t T) bool {
	s.sync.Lock()
	defer s.sync.Unlock()
	_, tPresent := s.internal[t]
	return tPresent
}

func (s *MapSet[T]) ContainAll(ts ...T) bool {
	s.sync.Lock()
	defer s.sync.Unlock()
	for _, t := range ts {
		if _, tPresent := s.internal[t]; !tPresent {
			return false
		}
	}
	return true
}

func (s *MapSet[T]) Remove(t T) {
	s.sync.Lock()
	defer s.sync.Unlock()
	delete(s.internal, t)
}

func (s *MapSet[T]) RemoveAll(ts ...T) {
	s.sync.Lock()
	defer s.sync.Unlock()
	for _, t := range ts {
		delete(s.internal, t)
	}
}
