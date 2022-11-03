package binn

import (
	"fmt"
	"sync"
)

type Keeper interface {
	Add(*Bottle) error
	Get() (*Bottle, error)
}

type BottleStorage struct {
	size    int
	bottles []*Bottle
	mux     *sync.Mutex
}

func NewBottleStorage(size int) *BottleStorage {
	return &BottleStorage{
		size:    size,
		bottles: []*Bottle{},
		mux:     &sync.Mutex{},
	}
}

func (s *BottleStorage) Add(b *Bottle) error {
	s.mux.Lock()
	if len(s.bottles) >= s.size {
		return fmt.Errorf("storage is full")
	}
	s.bottles = append(s.bottles, b)
	s.mux.Unlock()
	return nil
}

func (s *BottleStorage) Get() (*Bottle, error) {
	s.mux.Lock()
	if len(s.bottles) == 0 {
		s.mux.Unlock()
		return nil, fmt.Errorf("storage has no containers")
	}
	b := s.bottles[0]
	s.bottles = s.bottles[1:]
	s.mux.Unlock()
	return b, nil
}
