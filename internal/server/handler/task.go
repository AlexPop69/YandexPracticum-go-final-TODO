package handler

import (
	"YandexPracticum-go-final-TODO/internal/storage"
	"YandexPracticum-go-final-TODO/internal/task"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

func AddTask(db *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received reqest POST AddTask")
		var err error
		var buf bytes.Buffer
		var t task.Task

		_, err = buf.ReadFrom(r.Body)
		if err != nil {
			log.Printf("Can`t read request body: %v\n", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err = json.Unmarshal(buf.Bytes(), &t); err != nil {
			log.Printf("Can`t unmarshal task: %v\n", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = task.Check(t)
		if err != nil {
			log.Println(err)
			json.NewEncoder(w).Encode(map[string]string{"error": string(err.Error())})
			return
		}

		date, _ := time.Parse("20060102", t.Date)
		if date.Before(time.Now()) || t.Repeat != "" {
			t.Date, _ = task.NextDate(time.Now(), t.Date, t.Repeat)
		}

		if t.Repeat == "d 1" {
			t.Date = time.Now().Format("20060102")
		}

		id, err := storage.Task.Add(db, &t)
		if err != nil {
			log.Fatalf("can't add task: %v", err)
			return
		}

		result := map[string]string{
			"id": strconv.Itoa(id),
		}

		resp, err := json.Marshal(result)
		if err != nil {
			log.Printf("Can`t marshal id: %v\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		_, err = w.Write(resp)
		if err != nil {
			log.Fatalf("can't write response: %v", err)
		}
		log.Printf("Task %s id:%v added successfully", t.Title, id)
	}

}
