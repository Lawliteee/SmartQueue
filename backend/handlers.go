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

	// Генерируем уникальный ID
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
		Participants:        []string{},
	}

	h.store.Save(queue)

	// Формируем ссылку
	link := "http://localhost:5173/queue/" + id // или твой порт фронта

	resp := QueueResponse{
		ID:   id,
		Name: req.Name,
		Link: link,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)

	// Логируем ссылку в консоль
	println("Новая очередь создана:", link)
}

// GET /api/queues/{id} – получение информации об очереди
func (h *Handler) GetQueue(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	queue, ok := h.store.Get(id)
	if !ok {
		http.Error(w, "Queue not found", http.StatusNotFound)
		return
	}

	// Для фронта отдаём только нужные поля
	type QueueInfo struct {
		Name         string `json:"name"`
		Participants int    `json:"participants"` // пока 0, потом добавим логику
		WaitTime     int    `json:"waitTime"`     // можно вычислять
		PeopleAhead  int    `json:"peopleAhead"`
	}
	info := QueueInfo{
		Name:         queue.Name,
		Participants: len(queue.Participants),
		WaitTime:     0,
		PeopleAhead:  0,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(info)
}
