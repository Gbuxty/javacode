package domain

import "errors"

var (
	ErrInsufficientFunds = errors.New("insufficient funds")
	ErrInvalidOperation  = errors.New("invalid operation type")
	ErrWalletNotFound    = errors.New("wallet not found")
	Errdb                = errors.New("db error")
	ErrAmout             = errors.New("amount cannot be 0")
)
