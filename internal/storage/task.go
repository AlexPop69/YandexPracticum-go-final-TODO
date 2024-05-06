package storage

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"YandexPracticum-go-final-TODO/internal/task"
)

const taskLimit = 10

func (s *Storage) Add(t *task.Task) (int, error) {
	ins, err := s.db.Exec(
		"INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)",
		sql.Named("date", t.Date),
		sql.Named("title", t.Title),
		sql.Named("comment", t.Comment),
		sql.Named("repeat", t.Repeat),
	)
	if err != nil {
		return 0, fmt.Errorf("can't add task: %v", err)
	}

	id, err := ins.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("can't het last insert id: %v", err)
	}

	return int(id), nil
}

func (s *Storage) GetList() ([]task.Task, error) {
	rows, err := s.db.Query(
		`SELECT * 
		 FROM scheduler
		 ORDER BY date ASC
		 LIMIT :limit`,
		sql.Named("limit", taskLimit),
	)
	if err != nil {
		log.Printf("can't get tasks by GetList: %v", err)
		return nil, fmt.Errorf("can't get tasks: %v", err)
	}

	defer rows.Close()

	var tasks []task.Task

	for rows.Next() {
		t := task.Task{}

		err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			log.Printf("can't get tasks by GetList: %v", err)
			return nil, fmt.Errorf("can't get tasks %v", err)
		}

		tasks = append(tasks, t)
	}

	return tasks, nil
}

func (s *Storage) SearchTasks(search string) ([]task.Task, error) {
	log.Printf("looking for a task with search parameter %s", search)

	var rows *sql.Rows
	var err error

	date, err := time.Parse("02.01.2006", search)
	if err != nil {
		log.Println("Search by text")

		rows, err = s.db.Query(
			`SELECT * 
			FROM scheduler
			WHERE title LIKE :target OR comment LIKE :target
			ORDER BY date LIMIT :limit`,
			sql.Named("target", "%"+search+"%"),
			sql.Named("limit", taskLimit),
		)
	} else {
		log.Println("Search by date")

		target := date.Format("20060102")

		rows, err = s.db.Query(
			`SELECT * 
			   FROM scheduler
			   WHERE date LIKE :target
			   ORDER BY date LIMIT :limit`,
			sql.Named("target", "%"+target+"%"),
			sql.Named("limit", taskLimit),
		)
	}

	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	defer rows.Close()

	var tasks []task.Task

	for rows.Next() {
		t := task.Task{}
		err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			return nil, fmt.Errorf("can't find tasks by SearchTasks %v", err)
		}

		tasks = append(tasks, t)
	}
	log.Printf("Found %d tasks", len(tasks))

	return tasks, nil
}
