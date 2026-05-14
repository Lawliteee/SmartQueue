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
	swaps *SwapStore
}

func NewHandler(store *Store, hub *Hub) *Handler {
	return &Handler{store: store, hub: hub, swaps: NewSwapStore()}
}

// Обслуживает WebSocket-подключения к очереди
func (h *Handler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	queueID := vars["id"]
	participantID := r.URL.Query().Get("participantId")

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

	h.hub.Subscribe(queueID, participantID, conn)
	defer h.hub.Unsubscribe(queueID, participantID, conn)
	
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
		ImFreeFeature: 		 req.ImFreeFeature,
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

	// текущий участник
	var currentParticipant *Participant
	var cpID, cpName sql.NullString

	h.store.db.QueryRow(`
		SELECT current_participant_id, current_participant_name
		FROM queues WHERE id = $1
	`, id).Scan(&cpID, &cpName)

	if cpID.Valid && cpID.String != "" {
		currentParticipant = &Participant{
			ID:   cpID.String,
			Name: cpName.String,
		}
	}

	// Считаем среднее время
	avgWait := float64(5)

	rows, err := h.store.db.Query(`
		SELECT wait_time FROM queue_wait_stats
		WHERE queue_id = $1
	`, id)

	if err == nil {
		defer rows.Close()

		total := 0
		count := 0

		for rows.Next() {
			var wt int
			if err := rows.Scan(&wt); err == nil {
				total += wt
				count++
			}
		}

		if count > 0 {
			avgWait = float64(total) / float64(count)
		}
	}

	// ETA для каждого
	etAs := make([]int, len(queue.Participants))

	for i := range queue.Participants {
		etAs[i] = int(avgWait * float64(i+1))
	}

	// текущее ожидание
	waitTime := 0
	if len(etAs) > 0 {
		waitTime = etAs[0]
	}

	// Ответ
	type QueueInfo struct {
		ID                 string        `json:"id"`
		Name               string        `json:"name"`
		Description        string        `json:"description"`
		StartTime          string        `json:"startTime"`
		CurrentNumber      int           `json:"currentNumber"`
		Participants       []Participant `json:"participants"`
		PeopleAhead        int           `json:"peopleAhead"`
		WaitTime           int           `json:"waitTime"`
		ETAs               []int         `json:"etAs"`
		Finished           bool          `json:"finished"`
		CurrentParticipant *Participant  `json:"currentParticipant"`
		MaxParticipants    int           `json:"maxParticipants"`
		ImFreeFeature 	   bool          `json:"imFreeFeature"`
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(QueueInfo{
		ID:                 queue.ID,
		Name:               queue.Name,
		Description:        queue.Description,
		StartTime:          queue.StartTime,
		CurrentNumber:      queue.CurrentNumber,
		Participants:       queue.Participants,
		PeopleAhead:        len(queue.Participants),
		WaitTime:           waitTime,
		ETAs:               etAs,
		Finished:           queue.Finished,
		CurrentParticipant: currentParticipant,
		MaxParticipants:    queue.MaxParticipants,
		ImFreeFeature: 		queue.ImFreeFeature,
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

	// Получаем обновлённую очередь
	updatedQueue, _ := h.store.Get(id)

	// Отправляем обновление всем через WebSocket
	h.hub.Broadcast(id, map[string]interface{}{
		"type": "queue_update",
		"data": updatedQueue,
	})

	type NextResponse struct {
		CurrentNumber      int           `json:"currentNumber"`
		Participants       []Participant `json:"participants"`
		CurrentParticipant *Participant  `json:"currentParticipant"`
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(NextResponse{
		CurrentNumber:      updatedQueue.CurrentNumber,
		Participants:       updatedQueue.Participants,
		CurrentParticipant: updatedQueue.CurrentParticipant,
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

// POST /api/queues/{id}/swap/request
func (h *Handler) SwapRequest(w http.ResponseWriter, r *http.Request) {
    queueID := mux.Vars(r)["id"]

    var req SwapRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid body", http.StatusBadRequest)
        return
    }

    offer := &SwapOffer{
        ID:       uuid.New().String(),
        QueueID:  queueID,
        FromID:   req.FromID,
        FromName: req.FromName,
        ToID:     req.ToID,
    }
    h.swaps.Add(offer)
	
	// Находим позицию fromID в очереди
	queue, _ := h.store.Get(queueID)
	fromPos := 0
	for i, p := range queue.Participants {
    if p.ID == req.FromID {
			fromPos = i + 1
			break
		}
	}


    // Шлём уведомление только целевому участнику через WS
    h.hub.BroadcastTo(queueID, req.ToID, map[string]interface{}{
        "type": "swap_request",
        "swapId":   offer.ID,
        "fromId":   req.FromID,
        "fromName": req.FromName,
		"fromPos":  fromPos,
    })

    w.WriteHeader(http.StatusOK)
}

// POST /api/queues/{id}/swap/respond
func (h *Handler) SwapRespond(w http.ResponseWriter, r *http.Request) {
    queueID := mux.Vars(r)["id"]

    var req SwapRespondRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid body", http.StatusBadRequest)
        return
    }

    offer, ok := h.swaps.Take(req.SwapID)
    if !ok {
        http.Error(w, "Swap not found or expired", http.StatusNotFound)
        return
    }

    if req.Accept {
        if err := h.store.SwapParticipants(queueID, offer.FromID, offer.ToID); err != nil {
            http.Error(w, "Swap failed", http.StatusInternalServerError)
            return
        }
        // Всем обновление очереди
        if q, ok := h.store.Get(queueID); ok {
            h.hub.Broadcast(queueID, map[string]interface{}{
                "type": "queue_update",
                "data": q,
            })
        }
    } else {
        // Достаём имя отказавшего участника
        var declinerName string
        h.store.db.QueryRow(`SELECT name FROM participants WHERE id = $1`, offer.ToID).Scan(&declinerName)

        h.hub.BroadcastTo(queueID, offer.FromID, map[string]interface{}{
            "type":     "swap_declined",
            "fromName": declinerName,
        })
    }

    w.WriteHeader(http.StatusOK)
}

// POST /api/queues/{id}/im-free – участник освободился, вызываем следующего
func (h *Handler) ImFree(w http.ResponseWriter, r *http.Request) {
    queueID := mux.Vars(r)["id"]

    queue, ok := h.store.Get(queueID)
    if !ok {
        http.Error(w, "Queue not found", http.StatusNotFound)
        return
    }

    // Если есть следующие — вызываем
    if len(queue.Participants) > 0 {
        h.store.ShiftParticipant(queueID)
    } else {
        // Очередь пуста — просто сбрасываем текущего
        h.store.ClearCurrentParticipant(queueID)
    }

    if q, ok := h.store.Get(queueID); ok {
        h.hub.Broadcast(queueID, map[string]interface{}{
            "type": "queue_update",
            "data": q,
        })
    }

    w.WriteHeader(http.StatusOK)
}

// POST /api/queues/{id}/skip
func (h *Handler) SkipMe(w http.ResponseWriter, r *http.Request) {
    queueID := mux.Vars(r)["id"]

    var body struct {
        ParticipantID string `json:"participantId"`
    }
    if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.ParticipantID == "" {
        http.Error(w, "Invalid body", http.StatusBadRequest)
        return
    }

    if err := h.store.SkipParticipant(queueID, body.ParticipantID); err != nil {
        http.Error(w, "Skip failed", http.StatusInternalServerError)
        return
    }

    h.store.ExpireSkips(queueID)

    if q, ok := h.store.Get(queueID); ok {
        h.hub.Broadcast(queueID, map[string]interface{}{
            "type": "queue_update",
            "data": q,
        })
    }
    w.WriteHeader(http.StatusOK)
}

// POST /api/queues/{id}/return
func (h *Handler) ReturnMe(w http.ResponseWriter, r *http.Request) {
    queueID := mux.Vars(r)["id"]

    var body struct {
        ParticipantID string `json:"participantId"`
    }
    if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.ParticipantID == "" {
        http.Error(w, "Invalid body", http.StatusBadRequest)
        return
    }

    if err := h.store.ReturnParticipant(queueID, body.ParticipantID); err != nil {
        http.Error(w, "Return failed", http.StatusInternalServerError)
        return
    }

    if q, ok := h.store.Get(queueID); ok {
        h.hub.Broadcast(queueID, map[string]interface{}{
            "type": "queue_update",
            "data": q,
        })
    }
    w.WriteHeader(http.StatusOK)
}
