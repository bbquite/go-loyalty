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
		SELECT id 
		FROM account 
		WHERE username = $1 
		LIMIT 1
	`

	row := storage.Db.QueryRowContext(storage.ctx, sqlString, username)
	err := row.Scan(&account.Id)
	if err != nil {
		return account, err
	}

	return account, nil
}

func (storage *DBStorage) GetAccountByLoginData(username string, password string) (models.Account, error) {
	var account models.Account

	sqlString := `
		SELECT id 
		FROM account 
		WHERE username = $1 AND password = $2
		LIMIT 1
	`

	row := storage.Db.QueryRowContext(storage.ctx, sqlString, username, password)
	err := row.Scan(&account.Id)
	if err != nil {
		return account, err
	}

	return account, nil
}

func (storage *DBStorage) CreateAccount(username string, password string) (uint32, error) {
	sqlString := `
		INSERT INTO account (username, password) 
		VALUES ($1, $2)
		RETURNING id
	`
	// pgx не поддерживает LastInsertId
	// result, err := storage.Db.ExecContext(storage.ctx, sqlString, username, password)

	var userID uint32
	row := storage.Db.QueryRowContext(storage.ctx, sqlString, username, password)
	row.Scan(&userID)

	if userID == 0 {
		return 0, errors.New("unspecified error while creating record")
	}

	return userID, nil
}

func (storage *DBStorage) CreatePurchase(accountID uint32, purchaseID string) error {
	sqlString := `
		INSERT INTO purchase (account_id, purchase_num, purchase_status) 
		VALUES ($1, $2, $3)
		RETURNING id
	`
	_, err := storage.Db.ExecContext(storage.ctx, sqlString, accountID, purchaseID, "NEW")
	if err != nil {
		return err
	}

	return nil
}

func (storage *DBStorage) GetPurchasesList(accountID uint32) ([]models.Purchase, error) {

	var result []models.Purchase

	sqlString := `
		SELECT purchase_num, purchase_status, uploaded_at
		FROM purchase 
		WHERE account_id = $1
	`

	rows, err := storage.Db.QueryContext(storage.ctx, sqlString, accountID)

	for rows.Next() {
		var purchaseItem models.Purchase

		err = rows.Scan(&purchaseItem.PurchaseNum, &purchaseItem.PurchaseStatus, &purchaseItem.UploadedAt)
		if err != nil {
			return nil, err
		}

		result = append(result, purchaseItem)
	}
	if err != nil {
		return nil, err
	}

	return result, nil
}
