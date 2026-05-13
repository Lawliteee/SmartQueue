package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	db := NewDB()
	defer db.Close()

	RunMigrations(db)

	store := NewStore(db)
	hub := NewHub()                   // создаём хаб
	handler := NewHandler(store, hub) // передаём хаб

	r := mux.NewRouter()
	api := r.PathPrefix("/api").Subrouter()

	// Аутентификация (публичные)
	api.HandleFunc("/auth/register", handler.Register).Methods("POST")
	api.HandleFunc("/auth/login", handler.Login).Methods("POST")

	// Очереди (публичные)
	api.HandleFunc("/queues", handler.CreateQueue).Methods("POST")
	api.HandleFunc("/queues/{id}", handler.GetQueue).Methods("GET")
	api.HandleFunc("/queues/{id}/join", handler.JoinQueue).Methods("POST")
	api.HandleFunc("/queues/{id}/participants/{participantId}", handler.LeaveQueue).Methods("DELETE")
	api.HandleFunc("/queues/{id}/swap/request", handler.SwapRequest).Methods("POST")
	api.HandleFunc("/queues/{id}/swap/respond", handler.SwapRespond).Methods("POST")
	api.HandleFunc("/queues/{id}/im-free", handler.ImFree).Methods("POST")

	// WebSocket для live обновлений
	api.HandleFunc("/ws/queue/{id}", handler.HandleWebSocket).Methods("GET")
	// (только для авторизованных)
	admin := api.PathPrefix("/admin").Subrouter()
	admin.Use(AuthMiddleware)
	admin.HandleFunc("/queues/{id}/next", handler.CallNext).Methods("POST")
	admin.HandleFunc("/queues/{id}/finish", handler.FinishQueue).Methods("POST")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	port := ":8080"
	log.Printf("Сервер запущен на http://localhost%s", port)
	log.Fatal(http.ListenAndServe(port, c.Handler(r)))
}
