package domain

import (
	"github.com/google/uuid"
)

type OperationType string

const (
	Deposit  OperationType = "DEPOSIT"
	Withdraw OperationType = "WITHDRAW"
)

type Wallet struct {
	ID      uuid.UUID
	Balance int64
}

type WalletOperation struct {
	WalletID      uuid.UUID
	OperationType OperationType
	Amount        int64
}
