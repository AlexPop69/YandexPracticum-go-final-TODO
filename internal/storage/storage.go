package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

const (
	varEnvDBFile = "TODO_DBFILE"
	stdDbPath    = "./db"
	stdDbName    = "scheduler.db"
	dbDriver     = "sqlite"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath func() string) (*Storage, error) {
	dbPath := Path()

	if isExist(dbPath) {
		createDB(dbPath)
	} else {
		log.Println("Database is already exist")
	}

	db, err := sql.Open(dbDriver, dbPath)
	if err != nil {
		return nil, fmt.Errorf("can't open data base %w", err)
	}

	createNewTable(db)

	return &Storage{db: db}, nil
}

func Path() string {
	storagePath, exists := os.LookupEnv(varEnvDBFile)

	if !exists || storagePath == "" {
		storagePath = filepath.Join(stdDbPath, stdDbName)
		log.Printf(`Database storage address: %s`, storagePath)
	} else {
		log.Printf(`Database storage address %s retrieved from env variable "%s" `,
			storagePath,
			varEnvDBFile)
	}

	return storagePath
}

func isExist(dbPath string) bool {
	_, err := os.Stat(dbPath)

	return errors.Is(err, os.ErrNotExist)
}

func createDB(dbPath string) error {
	_, err := os.Create(dbPath)
	if err != nil {
		return fmt.Errorf("can't create storage by func createDB: %w", err)
	}

	log.Println("The database file has been created")

	return nil
}

func createNewTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS scheduler (
		   id INTEGER PRIMARY KEY AUTOINCREMENT,
		   date VARCHAR(8) NOT NULL,
		   title TEXT NOT NULL,
		   comment TEXT DEFAULT "",
		   repeat VARCHAR(128) NOT NULL
   		);
	
   		CREATE INDEX scheduler_date ON scheduler (date);
   `)

	if err != nil {
		return fmt.Errorf("can't create new table by func createNewTable: %w", err)
	}

	return nil
}
