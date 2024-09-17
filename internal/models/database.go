package models

import "time"

type Account struct {
	Id        uint32    `json:"id"`
	Username  string    `json:"login"`
	Password  string    `json:"password"`
	Balance   uint32    `json:"balance"`
	CreatedOn time.Time `json:"created_on"`
}

type OrderStatus struct {
	NEW        string `json:"NEW"`
	PROCESSING string `json:"PROCESSING"`
	INVALID    string `json:"INVALID"`
	PROCESSED  string `json:"PROCESSED"`
}

type Purchase struct {
	Id          uint32      `json:"id"`
	AccountID   uint32      `json:"account_id"`
	OrderNUM    uint32      `json:"order_num"`
	OrderStatus OrderStatus `json:"order_status"`
	UploadedAt  time.Time   `json:"uploaded_at"`
}

type TransactionType struct {
	IN  string `json:"IN"`
	OUT string `json:"OUT"`
}

type BalanceHistory struct {
	Id              uint32          `json:"id"`
	AccountID       uint32          `json:"account_id"`
	OrderID         uint32          `json:"order_id"`
	Amount          uint32          `json:"amount"`
	TransactionType TransactionType `json:"transaction_type"`
	ProcessedAt     time.Time       `json:"processed_at"`
}
