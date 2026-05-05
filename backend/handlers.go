package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Handler struct {
	store *Store
	hub   *Hub // WebSocket хаб
}

func NewHandler(store *Store, hub *Hub) *Handler {
	return &Handler{store: store, hub: hub}
}

// HandleWebSocket – обслуживает WebSocket-подключения к очереди
func (h *Handler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	queueID := vars["id"]

	// Проверяем существование очереди
	if _, ok := h.store.Get(queueID); !ok {
		http.Error(w, "Queue not found", http.StatusNotFound)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer conn.Close()

	h.hub.Subscribe(queueID, conn)
	defer h.hub.Unsubscribe(queueID, conn)

	// Ожидаем закрытия соединения (пока не обрабатываем входящие сообщения)
	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			break
		}
	}
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

	// Читаем текущего участника
	var currentParticipant *Participant
	var cpID, cpName sql.NullString
	h.store.db.QueryRow(`
		SELECT current_participant_id, current_participant_name
		FROM queues WHERE id = $1`, id,
	).Scan(&cpID, &cpName)
	if cpID.Valid && cpID.String != "" {
		currentParticipant = &Participant{ID: cpID.String, Name: cpName.String}
	}



	rows, err := h.store.db.Query(`
		SELECT wait_time FROM queue_wait_stats
		WHERE queue_id = $1
		`, id)

	avgWait := 5

	if err == nil {
		defer rows.Close()

		total := 0
		count := 0

		for rows.Next() {
			var wt int
			rows.Scan(&wt)
			total += wt
			count++
		}

		if count > 0 {
			avgWait = total / count
		}
	}



	type QueueInfo struct {
		ID                 string        `json:"id"`
		Name               string        `json:"name"`
		StartTime          string        `json:"startTime"`
		CurrentNumber      int           `json:"currentNumber"`
		Participants       []Participant `json:"participants"`
		PeopleAhead        int           `json:"peopleAhead"`
		WaitTime           int           `json:"waitTime"`
		Finished           bool          `json:"finished"`
		CurrentParticipant *Participant  `json:"currentParticipant"`
		MaxParticipants    int           `json:"maxParticipants"`
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(QueueInfo{
		ID:                 queue.ID,
		Name:               queue.Name,
		StartTime:          queue.StartTime,
		CurrentNumber:      queue.CurrentNumber,
		Participants:       queue.Participants,
		PeopleAhead:        len(queue.Participants),
		WaitTime:           avgWait * len(queue.Participants), // новая логика среднего времени ожидания
		Finished:           queue.Finished,
		CurrentParticipant: currentParticipant,
		MaxParticipants:    queue.MaxParticipants,
	})
}

// POST /api/queues/{id}/join – вступление в очередь
func (h *Handler) JoinQueue(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	queue, ok := h.store.Get(id)
	if !ok {
		http.Error(w, "Queue not found", http.StatusNotFound)
		return
	}

	var req JoinQueueRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if queue.MaxParticipants > 0 && len(queue.Participants) >= queue.MaxParticipants {
		http.Error(w, "Queue is full", http.StatusConflict)
		return
	}

	participant := Participant{
		ID:       uuid.New().String(),
		Name:     req.Name,
		Priority: 0,
	}

	h.store.AddParticipant(id, participant)

	// Отправляем обновление всем, кто слушает WebSocket
	if q, ok := h.store.Get(id); ok {
		h.hub.Broadcast(id, map[string]interface{}{
			"type": "queue_update",
			"data": q,
		})
	}

	queue, _ = h.store.Get(id)

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

	// Отправляем обновление всем, кто слушает WebSocket
	if q, ok := h.store.Get(id); ok {
		h.hub.Broadcast(id, map[string]interface{}{
			"type": "queue_update",
			"data": q,
		})
	}

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

	// Отправляем обновление всем, кто слушает WebSocket
	if q, ok := h.store.Get(id); ok {
		h.hub.Broadcast(id, map[string]interface{}{
			"type": "queue_update",
			"data": q,
		})
	}
	w.WriteHeader(http.StatusOK)
}

// DELETE /api/queues/{id}/participants/{participantId} – покинуть очередь
func (h *Handler) LeaveQueue(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	h.store.RemoveParticipant(vars["id"], vars["participantId"])

	// Отправляем обновление всем, кто слушает WebSocket
	if q, ok := h.store.Get(vars["id"]); ok {
		h.hub.Broadcast(vars["id"], map[string]interface{}{
			"type": "queue_update",
			"data": q,
		})
	}
	w.WriteHeader(http.StatusOK)
}
