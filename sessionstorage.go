package main

import "sync"

type SessionStorage struct {
	mux     *sync.Mutex
	storage map[int64]*RequestDTO
	arr     []int64
}

func NewSessionStorage() *SessionStorage {
	return &SessionStorage{
		mux:     &sync.Mutex{},
		storage: make(map[int64]*RequestDTO),
		arr:     make([]int64, 0),
	}
}

func (s *SessionStorage) Store(key int64, value *RequestDTO) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.arr = append(s.arr, key)
	if len(s.arr) > 100 {
		for _, k := range s.arr[:50] {
			delete(s.storage, k)
		}
		s.arr = s.arr[50:]
	}
	s.storage[key] = value
}

func (s *SessionStorage) Load(key int64) (*RequestDTO, bool) {
	s.mux.Lock()
	value, ok := s.storage[key]
	s.mux.Unlock()
	return value, ok
}
