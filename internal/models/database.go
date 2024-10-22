package models

import "time"

type Account struct {
	Id        uint32    `json:"id"`
	Username  string    `json:"login"`
	Password  string    `json:"password"`
	Balance   uint32    `json:"balance"`
	CreatedOn time.Time `json:"created_on"`
}

// type OrderStatus struct {
// 	NEW        string `json:"NEW"`
// 	PROCESSING string `json:"PROCESSING"`
// 	INVALID    string `json:"INVALID"`
// 	PROCESSED  string `json:"PROCESSED"`
// }

type Purchase struct {
	Id             uint32    `json:"id,omitempty"`
	AccountID      uint32    `json:"account_id,omitempty"`
	PurchaseNum    string    `json:"number"`
	PurchaseStatus string    `json:"status"`
	UploadedAt     time.Time `json:"uploaded_at"`
}

type Balance struct {
	Id        uint32  `json:"id,omitempty"`
	AccountId uint32  `json:"account_id,omitempty"`
	Amount    float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

//type TransactionType struct {
//	IN  string `json:"IN"`
//	OUT string `json:"OUT"`
//}

type BalanceHistory struct {
	Id              uint32    `json:"id,omitempty"`
	AccountID       uint32    `json:"account_id,omitempty"`
	PurchaseID      uint32    `json:"order"`
	Amount          uint32    `json:"sum"`
	TransactionType string    `json:"transaction_type,omitempty"`
	ProcessedAt     time.Time `json:"processed_at"`
}
