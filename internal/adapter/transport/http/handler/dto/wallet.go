package dto

import (
	"fmt"
	"java-code/internal/domain"

	"github.com/google/uuid"
)

type WalletOperationRequest struct {
	WalletID      string `json:"walletId" binding:"required"`
	OperationType string `json:"operationType" binding:"required,oneof=DEPOSIT WITHDRAW"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
}

type WalletBalanceResponse struct {
	WalletID string `json:"walletId"`
	Balance  int64  `json:"balance"`
}

func ToDomainWalletOperation(r WalletOperationRequest) (domain.WalletOperation, error) {
	id, err := uuid.Parse(r.WalletID)
	if err != nil {
		return domain.WalletOperation{}, fmt.Errorf("invalid wallet ID %q: %w", r.WalletID, err)
	}

	return domain.WalletOperation{
		WalletID:      id,
		OperationType: domain.OperationType(r.OperationType),
		Amount:        r.Amount,
	}, nil
}

func ToDtoWalletBalanceRepsonse(r domain.Wallet) WalletBalanceResponse {
	return WalletBalanceResponse{
		WalletID: r.ID.String(),
		Balance:  r.Balance,
	}
}
