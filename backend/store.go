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

func (s *Store) AddParticipant(queueID string, p Participant) {
	s.mu.Lock()
	defer s.mu.Unlock()
	q, ok := s.queues[queueID]
	if !ok {
		return
	}
	q.Participants = append(q.Participants, p)
}

// ShiftParticipant удаляет первого участника и увеличивает счётчик
func (s *Store) ShiftParticipant(queueID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	q, ok := s.queues[queueID]
	if !ok || len(q.Participants) == 0 {
		return
	}
	q.Participants = q.Participants[1:]
	q.CurrentNumber++
}

func (s *Store) FinishQueue(queueID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	q, ok := s.queues[queueID]
	if !ok {
		return
	}
	q.Finished = true
}