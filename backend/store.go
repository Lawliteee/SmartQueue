package main

import (
  "sync"
)

type Store struct {
  queues map[string]*Queue
  mu     sync.RWMutex
}

func NewStore() *Store {
  return &Store{
    queues: make(map[string]*Queue),
  }
}

func (s *Store) Save(queue *Queue) {
  s.mu.Lock()
  defer s.mu.Unlock()
  s.queues[queue.ID] = queue
}

func (s *Store) Get(id string) (*Queue, bool) {
  s.mu.RLock()
  defer s.mu.RUnlock()
  q, ok := s.queues[id]
  return q, ok
}