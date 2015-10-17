package datatypes

import (
	"sync"
)

type SafeMap struct {
	sync.Mutex
	dict map[string]interface{}
}

func NewSafeMap() *SafeMap {
	return &SafeMap{
		dict: make(map[string]interface{}),
	}
}

func (s *SafeMap) Set(key string, value interface{}) {
	s.Lock()
	s.dict[key] = value
	s.Unlock()
}

func (s *SafeMap) Erase(key string) {
	s.Lock()
	delete(s.dict, key)
	s.Unlock()
}

func (s *SafeMap) Get(key string) (value interface{}, ok bool) {
	s.Lock()
	value, ok = s.dict[key]
	s.Unlock()
	return
}

func (s *SafeMap) Clone() (dict map[string]interface{}) {
	s.Lock()
	dict = s.dict
	s.Unlock()
	return
}

func (s *SafeMap) Len() (l int) {
	s.Lock()
	l = len(s.dict)
	s.Unlock()
	return
}
