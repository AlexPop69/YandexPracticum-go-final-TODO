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

	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	result, err := task.NextDate(now, date, repeat)
	if err != nil {
		log.Println(err)
	}

	realDate, _ := time.Parse("20060102", date)
	if repeat == "d 1" && realDate.After(now) {
		result = realDate.AddDate(0, 0, 1).Format("20060102")
	}

	if realDate == now && repeat == "" {
		result = ""
	}

	w.Header().Set("Content-Type", "string")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write([]byte(result))
	if err != nil {
		log.Println("Error write in func GetNextDate:", err)
	}
}
