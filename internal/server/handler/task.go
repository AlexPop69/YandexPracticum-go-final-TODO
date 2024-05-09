package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"YandexPracticum-go-final-TODO/internal/storage"
	"YandexPracticum-go-final-TODO/internal/task"
)

func AddTask(storage *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received reqest POST AddTask")

		var t task.Task

		err := json.NewDecoder(r.Body).Decode(&t)
		if err != nil {
			http.Error(w, `{"error":"Can't read request body"}`, http.StatusBadRequest)
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
			log.Fatalf("can't write response by GetTasks: %v", err)
		} else {
			log.Printf("GetTasks is successful. %d tasks found", len(tasks))
		}
	}
}

func GetTask(storage *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received reqest GET task")

		var err error

		id := r.URL.Query().Get("id")

		if id == "" {
			log.Println("ID is empty")
			json.NewEncoder(w).Encode(map[string]string{"error": "id is empty"})
			return
		}

		_, err = strconv.Atoi(id)
		if err != nil {
			log.Println("incorrect ID, id is not number")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "incorrect task ID"})
			return
		}

		task, err := storage.GetTask(id)
		if err != nil {
			log.Printf("can't get task: %v", err)
			http.Error(w, err.Error(), http.StatusNoContent)
			json.NewEncoder(w).Encode(map[string]string{"error": "can't get task"})
			return
		}

		resp, err := json.Marshal(task)
		if err != nil {
			log.Printf("Can`t marshal task: %v\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": string(err.Error())})
			return
		}

		w.Header().Set("Content-Type", "application/json")

		_, err = w.Write(resp)
		if err != nil {
			log.Fatalf("can't write response by GetTask: %v", err)
		} else {
			log.Println("GetTask is successful")
		}
	}
}

func UpdateTask(storage *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received reqest UpdateTask")

		var t task.Task

		err := json.NewDecoder(r.Body).Decode(&t)
		if err != nil {
			http.Error(w, `{"error":"Can't read request body"}`, http.StatusBadRequest)
			return
		}

		err = task.Check(t)
		if err != nil {
			log.Println(err)
			json.NewEncoder(w).Encode(map[string]string{"error": string(err.Error())})
			return
		}

		err = storage.Update(t)
		if err != nil {
			log.Println(err)
			json.NewEncoder(w).Encode(map[string]string{"error": string(err.Error())})
			return
		}

		t, err = storage.GetTask(t.ID)
		if err != nil {
			log.Println(err)
			json.NewEncoder(w).Encode(map[string]string{"error": string(err.Error())})
			return
		}

		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(map[string]string{})
	}
}
