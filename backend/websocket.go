package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Управляет всеми вебсокет соединениями по очереди.
type Hub struct {
	// rooms: queueID - множество активных соединений
	rooms map[string]map[*websocket.Conn]bool
	clients map[string]map[string]*websocket.Conn
	history  map[string][]ChatMessage
	mu    sync.RWMutex
}

// Создаёт новый Hub.
func NewHub() *Hub {
	return &Hub{
		rooms: make(map[string]map[*websocket.Conn]bool),
		clients: make(map[string]map[string]*websocket.Conn),
		history: make(map[string][]ChatMessage),
	}
}

// Добавляет соединение в комнату очереди.
func (h *Hub) Subscribe(queueID, participantID string, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.rooms[queueID]; !ok {
		h.rooms[queueID] = make(map[*websocket.Conn]bool)
		h.clients[queueID] = make(map[string]*websocket.Conn)
	}
	h.rooms[queueID][conn] = true
	if participantID != "" {
		h.clients[queueID][participantID] = conn
	}
	// Отправляем историю новому подключению
    if msgs, ok := h.history[queueID]; ok && len(msgs) > 0 {
        data, err := json.Marshal(map[string]interface{}{
            "type": "chat_history",
            "messages": msgs,
        })
        if err == nil {
            conn.WriteMessage(websocket.TextMessage, data)
        }
    }
	log.Printf("WebSocket подключён к очереди %s (всего соединений: %d)", queueID, len(h.rooms[queueID]))
}

// Удаляет соединение из комнаты.
func (h *Hub) Unsubscribe(queueID, participantID string, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if clients, ok := h.rooms[queueID]; ok {
		delete(clients, conn)
		if len(clients) == 0 {
			delete(h.rooms, queueID)
		}
	}
	if participantID != "" && h.clients[queueID] != nil {
		delete(h.clients[queueID], participantID)
	}
	log.Printf("WebSocket отключён от очереди %s", queueID)
}

// Отправляет JSON сообщение всем клиентам в комнате очереди.
func (h *Hub) Broadcast(queueID string, message interface{}) {
	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("Broadcast marshal error: %v", err)
		return
	}
	h.mu.RLock()
	defer h.mu.RUnlock()
	for conn := range h.rooms[queueID] {
		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			log.Printf("WebSocket write error: %v", err)
			conn.Close()
		}
	}
}

// Отправляет JSON-сообщение конкретному участнику
func (h *Hub) BroadcastTo(queueID, participantID string, message interface{}) {
	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("BroadcastTo marshal error: %v", err)
		return
	}
	h.mu.RLock()
	defer h.mu.RUnlock()
	if conn, ok := h.clients[queueID][participantID]; ok {
		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			log.Printf("BroadcastTo write error: %v", err)
			conn.Close()
		}
	}
}

type ChatMessage struct {
    SenderID   string `json:"senderId"`
    SenderName string `json:"senderName"`
    Text       string `json:"text"`
}

func (h *Hub) AddToHistory(queueID string, msg ChatMessage) {
    h.mu.Lock()
    defer h.mu.Unlock()
    h.history[queueID] = append(h.history[queueID], msg)
    if len(h.history[queueID]) > 20 {
        h.history[queueID] = h.history[queueID][len(h.history[queueID])-20:]
    }
}
