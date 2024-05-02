package main

import (
	"log"

	"YandexPracticum-go-final-TODO/internal/server"
	"YandexPracticum-go-final-TODO/internal/server/handler"
	"YandexPracticum-go-final-TODO/internal/storage"

	"github.com/go-chi/chi"
)

func main() {
	db, err := storage.New(storage.Path)
	if err != nil {
		log.Fatal("can't init storage", err)
	}
	_ = db

	r := chi.NewRouter()

	r.Handle("/*", handler.GetFront())

	server := new(server.Server)
	if err := server.Run(r); err != nil {
		log.Fatalf("Server can't start: %v", err)
		return
	}

	log.Println("Server stopped")
}
