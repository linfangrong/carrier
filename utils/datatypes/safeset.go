package datatypes

import (
	"sync"
)

type SafeSet struct {
	sync.Mutex
	dict map[string]bool
}

func NewSafeSet() *SafeSet {
	return &SafeSet{
		dict: make(map[string]bool),
	}
}

func (s *SafeSet) Add(element string) {
	s.Lock()
	defer s.Unlock()
	s.dict[element] = true
}

func (s *SafeSet) Erase(element string) {
	s.Lock()
	defer s.Unlock()
	delete(s.dict, element)
}

func (s *SafeSet) Contains(element string) bool {
	s.Lock()
	defer s.Unlock()
	_, ok := s.dict[element]
	return ok
}
