package storage

import (
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type DBStorage struct {
	db *sql.DB
	// mx sync.RWMutex
}

func NewDBStorage(databaseURI string) (*DBStorage, error) {
	db, err := sql.Open("pgx", databaseURI)
	if err != nil {
		return &DBStorage{}, err
	}

	return &DBStorage{
		db: db,
	}, nil
}

func (storage *DBStorage) RegisterUser() {
	log.Print("TEST")
}
