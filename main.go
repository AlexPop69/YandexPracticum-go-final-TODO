package main

import (
	"YandexPracticum-go-final-TODO/internal/storage"
	"YandexPracticum-go-final-TODO/server"
	"log"
	"net/http"
)

const (
	webDir = "./web"
)

func main() {
	db, err := storage.New(storage.Path)
	if err != nil {
		log.Fatal("can't init storage", err)
	}
	_ = db

	// получаем фронт
	http.Handle("/", http.FileServer(http.Dir(webDir)))
	log.Printf("Loaded frontEnd files from %s\n", webDir)

	server := new(server.Server)

	if err := server.Run(); err != nil {
		log.Fatalf("Server can't start: %v", err)
		return
	}

	log.Println("Server is stopped")
}
