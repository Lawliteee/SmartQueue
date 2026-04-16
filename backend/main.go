package main

import (
  "log"
  "net/http"

  "github.com/gorilla/mux"
  "github.com/rs/cors"
)

func main() {
  store := NewStore()
  handler := NewHandler(store)

  r := mux.NewRouter()

  // API маршруты
  api := r.PathPrefix("/api").Subrouter()
  api.HandleFunc("/queues", handler.CreateQueue).Methods("POST")
  api.HandleFunc("/queues/{id}", handler.GetQueue).Methods("GET")

  // Настройка CORS (чтобы фронт на другом порту мог обращаться)
  c := cors.New(cors.Options{
    AllowedOrigins:   []string{"http://localhost:5173"}, // порт Vite по умолчанию
    AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
    AllowedHeaders:   []string{"Content-Type"},
    AllowCredentials: true,
  })

  handlerWithCORS := c.Handler(r)

  port := ":8080"
  log.Printf("Сервер запущен на http://localhost%s", port)
  log.Fatal(http.ListenAndServe(port, handlerWithCORS))
}