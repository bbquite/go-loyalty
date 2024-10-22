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

	var userID uint32
	row := storage.Db.QueryRowContext(storage.ctx, sqlString, username, password)
	err := row.Scan(&userID)
	if err != nil {
		return 0, err
	}

	if userID == 0 {
		return 0, errors.New("unspecified error while creating record")
	}

	sqlStringBalance := `
		INSERT INTO balance (account_id) 
		VALUES ($1)
	`

	_, err = storage.Db.ExecContext(storage.ctx, sqlStringBalance, userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (storage *DBStorage) GetAccountBalance(accountID uint32) (models.Balance, error) {
	var balance models.Balance

	sqlString := `
		SELECT amount, withdrawn
		FROM balance 
		WHERE account_id = $1
		LIMIT 1
	`

	row := storage.Db.QueryRowContext(storage.ctx, sqlString, accountID)
	err := row.Scan(&balance.Amount, &balance.Withdrawn)
	if err != nil {
		return balance, err
	}

	return balance, nil
}

func (storage *DBStorage) GetAccountBalanceHistory(accountID uint32, trType string) ([]models.BalanceHistory, error) {

	var result []models.BalanceHistory

	sqlString := `
		SELECT purchase_id, amount, processed_at
		FROM balance_history 
		WHERE account_id = $1 AND transaction_type = $2
		ORDER BY processed_at ASC
	`

	rows, err := storage.Db.QueryContext(storage.ctx, sqlString, accountID, trType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var balanceItem models.BalanceHistory

		err = rows.Scan(&balanceItem.PurchaseID, &balanceItem.Amount, &balanceItem.ProcessedAt)
		if err != nil {
			return nil, err
		}

		result = append(result, balanceItem)
	}

	return result, nil
}

func (storage *DBStorage) GetPurchase(purchaseID string) (models.Purchase, error) {
	var purchase models.Purchase

	sqlString := `
		SELECT id, account_id, purchase_num, purchase_status, uploaded_at
		FROM purchase 
		WHERE purchase_num = $1
		LIMIT 1
	`

	row := storage.Db.QueryRowContext(storage.ctx, sqlString, purchaseID)
	err := row.Scan(&purchase.Id, &purchase.AccountID, &purchase.PurchaseNum, &purchase.PurchaseStatus, &purchase.UploadedAt)
	if err != nil {
		return purchase, err
	}

	return purchase, nil
}

func (storage *DBStorage) CreatePurchase(accountID uint32, purchaseID string) error {
	sqlString := `
		INSERT INTO purchase (account_id, purchase_num, purchase_status) 
		VALUES ($1, $2, $3)
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
		ORDER BY uploaded_at ASC
	`

	rows, err := storage.Db.QueryContext(storage.ctx, sqlString, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var purchaseItem models.Purchase

		err = rows.Scan(&purchaseItem.PurchaseNum, &purchaseItem.PurchaseStatus, &purchaseItem.UploadedAt)
		if err != nil {
			return nil, err
		}

		result = append(result, purchaseItem)
	}

	return result, nil
}
