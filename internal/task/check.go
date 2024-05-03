package task

import (
	"fmt"
	"log"
	"time"
)

func Check(t Task) error {
	log.Println("checking task")

	if t.Title == "" {
		return fmt.Errorf("title is empty")
	}

	if t.Date != "" {
		_, err := time.Parse("20060102", t.Date)
		if err != nil {
			return fmt.Errorf("date is in a format other than 20060102")
		}
	}

	_, err := NextDate(time.Now(), t.Date, t.Repeat)
	if err != nil {
		return err
	}

	return nil
}
