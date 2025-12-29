package models

import (
	"github.com/google/uuid"
)

type TransactionRequest struct {
	WalletID      string  `json:"walletId"`
	OperationType string  `json:"operationType"`
	Amount        float64 `json:"amount"`
}

type WalletResponse struct {
	ID      uuid.UUID `json:"id"`
	Status  string    `json:"status"`
	Balance float64   `json:"balance"`
}
