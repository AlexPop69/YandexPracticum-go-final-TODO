package main

import (
	"YandexPracticum-go-final-TODO/server"
	"log"
	"net/http"
)

const (
	WebDir = "./web"
)

func main() {

	// получаем фронт
	http.Handle("/", http.FileServer(http.Dir(WebDir)))

	server := new(server.Server)

	if err := server.Run(); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
		return
	}

	log.Println("Сервер остановлен")
}
