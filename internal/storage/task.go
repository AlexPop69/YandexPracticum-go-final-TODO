package storage

import (
	"database/sql"
	"errors"
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
		return 0, fmt.Errorf("can't get last insert id: %v", err)
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

func (s *Storage) GetTask(id string) (task.Task, error) {
	log.Println("Search task by ID")

	row := s.db.QueryRow(
		`SELECT * 
		FROM scheduler 
		WHERE id = :id`,
		sql.Named("id", id),
	)

	var t task.Task

	err := row.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if err != nil {
		log.Println("can't get task by id:", id, err)

		if errors.Is(err, sql.ErrNoRows) {
			return task.Task{}, errors.New(" ")
		}
		return task.Task{}, err
	}

	return t, nil
}

func (s *Storage) Update(t task.Task) error {
	log.Printf("Update task by ID:%s", t.ID)

	_, err := s.db.Exec(
		`UPDATE scheduler 
		SET date=:date, title= :title, comment= :comment, repeat= :repeat
		WHERE id= :id`,
		sql.Named("date", t.Date),
		sql.Named("title", t.Title),
		sql.Named("comment", t.Comment),
		sql.Named("repeat", t.Repeat),
		sql.Named("id", t.ID),
	)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("can't update task: %v", err)
	}

	log.Println("update successful")

	return nil
}

func (s *Storage) DoneTask(id string) error {
	log.Println("Done task ID:", id)

	t, err := s.GetTask(id)
	if err != nil {
		log.Println(err)
		return errors.New("task not found")
	}

	if t.Repeat == "" {
		log.Println("Repeat is empty, task will delete")
		err = s.DelTask(id)
		log.Println(err)

		return nil
	}

	if t.Repeat != "" {
		t.Date, err = task.NextDate(time.Now(), t.Date, t.Repeat)
		if err != nil {
			log.Println(err)
			return err
		}

		err = s.Update(t)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	log.Println("task is done id:", id)

	return nil
}

func (s *Storage) DelTask(id string) error {
	log.Println("Delete task ID:", id)

	_, err := s.db.Exec(
		`DELETE
		FROM scheduler
		WHERE id = :id`,
		sql.Named("id", id),
	)
	if err != nil {
		log.Println("can't delete task:", err)
		return errors.New("task not found")
	}

	log.Println("delete task successful")

	return nil
}
