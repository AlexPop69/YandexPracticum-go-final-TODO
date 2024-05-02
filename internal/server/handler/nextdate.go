package handler

import (
	"log"
	"net/http"
	"time"

	"YandexPracticum-go-final-TODO/internal/task"
)

func GetNextDate(w http.ResponseWriter, r *http.Request) {
	log.Println("Received reqest GetNextDate")

	r.ParseForm()

	now, err := time.Parse("20060102", r.FormValue("now"))
	if err != nil {
		log.Fatalf("Incorrect now date: %v", err)
	}

	result, err := task.NextDate(now,
		r.FormValue("date"),
		r.FormValue("repeat"),
	)
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "string")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write([]byte(result))
	if err != nil {
		log.Println("Error write in func GetNextDate:", err)
	}
}
