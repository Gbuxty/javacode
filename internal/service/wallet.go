package service

import (
	"context"
	"fmt"
	"java-code/internal/domain"

	"github.com/google/uuid"
)

type Wallet struct {
	repo      WalletRepository
	txManager TxManager
}

type WalletRepository interface {
	DepositBalance(ctx context.Context, op domain.WalletOperation) error
	WithdrawBalance(ctx context.Context, op domain.WalletOperation) error
	GetWalletBalance(ctx context.Context, id uuid.UUID) (domain.Wallet, error)
}

type TxManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}

func NewWallet(repo WalletRepository, txManager TxManager) *Wallet {
	return &Wallet{
		repo:      repo,
		txManager: txManager,
	}
}

func (s Wallet) OperationWithWallet(ctx context.Context, op domain.WalletOperation) error {
	switch op.OperationType {
	case domain.Deposit:
		if op.Amount <= 0 {
			return domain.ErrAmout
		}
		if err := s.repo.DepositBalance(ctx, op); err != nil {
			return fmt.Errorf("failed deposit balance:%w", err)
		}
	case domain.Withdraw:
		if err := s.txManager.Do(ctx, func(ctx context.Context) error {
			wallet, err := s.repo.GetWalletBalance(ctx, op.WalletID)
			if err != nil {
				return fmt.Errorf("failed get balance:%w", err)
			}
			if wallet.Balance < op.Amount {
				return domain.ErrInsufficientFunds
			}
			if err := s.repo.WithdrawBalance(ctx, op); err != nil {
				return fmt.Errorf("failed withdraw:%w", err)
			}

			return nil
		}); err != nil {
			return fmt.Errorf("txManager: %w", err)
		}
	default:
		return domain.ErrInvalidOperation
	}
	return nil
}
func (s Wallet) GetWalletByID(ctx context.Context, id uuid.UUID) (domain.Wallet, error) {
	


	wallet, err := s.repo.GetWalletBalance(ctx, id)
	if err != nil {
		return domain.Wallet{}, fmt.Errorf("failed get balance:%w,id:%s", err, id)
	}
	return wallet, nil
}
