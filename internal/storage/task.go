package storage

import (
	"YandexPracticum-go-final-TODO/internal/task"
	"database/sql"
	"fmt"
)

type Task interface {
	Add(t *task.Task) (int, error)
}

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
