package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Handler struct {
	store *Store
}

func NewHandler(store *Store) *Handler {
	return &Handler{store: store}
}

// POST /api/queues – создание очереди
func (h *Handler) CreateQueue(w http.ResponseWriter, r *http.Request) {
	var req CreateQueueRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	id := uuid.New().String()

	queue := &Queue{
		ID:                  id,
		Name:                req.Name,
		Description:         req.Description,
		StartTime:           req.StartTime,
		MaxParticipants:     req.MaxParticipants,
		HasPriority:         req.HasPriority,
		PriorityCount:       req.PriorityCount,
		InitialPriority:     req.InitialPriority,
		AnonymousChat:       req.AnonymousChat,
		SystemNotifications: req.SystemNotifications,
		SwapPositions:       req.SwapPositions,
		Admins:              req.Admins,
		CreatedAt:           time.Now(),
		Participants:        []Participant{},
		CurrentNumber:       0,
		Finished:            false,
	}

	if err := h.store.Save(queue); err != nil {
		http.Error(w, "Failed to create queue", http.StatusInternalServerError)
		return
	}

	link := "http://localhost:5173/queue/" + id

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(QueueResponse{
		ID:   id,
		Name: req.Name,
		Link: link,
	})
}

// GET /api/queues/{id} – информация об очереди
func (h *Handler) GetQueue(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	queue, ok := h.store.Get(id)
	if !ok {
		http.Error(w, "Queue not found", http.StatusNotFound)
		return
	}

	type QueueInfo struct {
		ID            string        `json:"id"`
		Name          string        `json:"name"`
		StartTime     string        `json:"startTime"`
		CurrentNumber int           `json:"currentNumber"`
		Participants  []Participant `json:"participants"`
		PeopleAhead   int           `json:"peopleAhead"`
		WaitTime      int           `json:"waitTime"`
		Finished      bool          `json:"finished"`
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(QueueInfo{
		ID:            queue.ID,
		Name:          queue.Name,
		StartTime:     queue.StartTime,
		CurrentNumber: queue.CurrentNumber,
		Participants:  queue.Participants,
		PeopleAhead:   len(queue.Participants),
		WaitTime:      len(queue.Participants) * 5, // заглушка: 5 мин на человека
		Finished:      queue.Finished,
	})
}

// POST /api/queues/{id}/join – вступление в очередь
func (h *Handler) JoinQueue(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if _, ok := h.store.Get(id); !ok {
		http.Error(w, "Queue not found", http.StatusNotFound)
		return
	}

	var req JoinQueueRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	participant := Participant{
		ID:       uuid.New().String(),
		Name:     req.Name,
		Priority: 0,
	}

	h.store.AddParticipant(id, participant)

	queue, _ := h.store.Get(id)

	type JoinResponse struct {
		ParticipantID string `json:"participantId"`
		Position      int    `json:"position"`
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(JoinResponse{
		ParticipantID: participant.ID,
		Position:      len(queue.Participants),
	})
}

// POST /api/admin/queues/{id}/next – вызвать следующего участника
func (h *Handler) CallNext(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	queue, ok := h.store.Get(id)
	if !ok {
		http.Error(w, "Queue not found", http.StatusNotFound)
		return
	}

	if len(queue.Participants) == 0 {
		http.Error(w, "No participants", http.StatusBadRequest)
		return
	}

	h.store.ShiftParticipant(id)
	queue, _ = h.store.Get(id)

	type NextResponse struct {
		CurrentNumber int           `json:"currentNumber"`
		Participants  []Participant `json:"participants"`
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(NextResponse{
		CurrentNumber: queue.CurrentNumber,
		Participants:  queue.Participants,
	})
}

// POST /api/admin/queues/{id}/finish – завершить очередь
func (h *Handler) FinishQueue(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if _, ok := h.store.Get(id); !ok {
		http.Error(w, "Queue not found", http.StatusNotFound)
		return
	}

	h.store.FinishQueue(id)
	w.WriteHeader(http.StatusOK)
}
