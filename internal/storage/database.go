package storage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/bbquite/go-loyalty/internal/models"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type DBStorage struct {
	Db  *sql.DB
	ctx context.Context
	// mx sync.RWMutex
}

func NewDBStorage(ctx context.Context, databaseURI string) (*DBStorage, error) {
	db, err := sql.Open("pgx", databaseURI)
	if err != nil {
		return &DBStorage{}, err
	}

	return &DBStorage{
		Db:  db,
		ctx: ctx,
	}, nil
}

func (storage *DBStorage) GetAccountByUsername(username string) (models.Account, error) {
	var account models.Account

	sqlString := `
		SELECT id, username 
		FROM account 
		WHERE username = $1 
		LIMIT 1
	`

	row := storage.Db.QueryRowContext(storage.ctx, sqlString, username)
	err := row.Scan(&account.Id, &account.Username)
	if err != nil {
		return account, err
	}

	return account, nil
}

func (storage *DBStorage) CreateAccount(username string, password string) (int64, error) {
	sqlString := `
		INSERT INTO account (username, password) 
		VALUES ($1, $2)
		RETURNING id
	`
	// pgx не поддерживает LastInsertId
	// result, err := storage.Db.ExecContext(storage.ctx, sqlString, username, password)

	var userID int64
	row := storage.Db.QueryRowContext(storage.ctx, sqlString, username, password)
	row.Scan(&userID)

	if userID == 0 {
		return 0, errors.New("unspecified error while creating record")
	}

	return userID, nil
}
