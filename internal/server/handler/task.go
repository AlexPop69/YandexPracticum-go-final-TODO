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

func AddTask(storage *storage.Storage) http.HandlerFunc {
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
		if date.Before(time.Now()) {
			t.Date, _ = task.NextDate(time.Now(), t.Date, t.Repeat)
		}

		id, err := storage.Add(&t)
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

func GetTasks(db *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received reqest GET Tasks")

		var tasks []task.Task
		var err error

		search := r.URL.Query().Get("search")

		if search == "" {
			tasks, err = db.GetList()
			if err != nil {
				log.Fatalf("can't get tasks: %v", err)
			}
		}

		if search != "" {
			tasks, err = db.SearchTasks(search)
			if err != nil {
				log.Fatalf("can't find tasks: %v", err)
			}
		}

		if len(tasks) == 0 {
			tasks = []task.Task{}
		}

		result := map[string][]task.Task{
			"tasks": tasks,
		}

		resp, err := json.Marshal(result)
		if err != nil {
			log.Printf("Can`t marshal tasks: %v\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": string(err.Error())})
			return
		}

		w.Header().Set("Content-Type", "application/json")

		_, err = w.Write(resp)
		if err != nil {
			log.Fatalf("can't write response: %v", err)
		} else {
			log.Printf("GetTasks is successful. %d tasks found", len(tasks))
		}
	}
}
