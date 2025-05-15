package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"java-code/internal/domain"

	"github.com/google/uuid"
)

type WalletRepository struct {
	executor *Executor
}

func NewWalletRepository(executor *Executor) *WalletRepository {
	return &WalletRepository{executor: executor}
}

func (w WalletRepository) DepositBalance(ctx context.Context, op domain.WalletOperation) error {
	executor := w.executor.Get(ctx)

	query := `
			INSERT INTO wallets (wallet_id, balance, created_at, updated_at)
			VALUES ($1, $2, NOW(), NOW())
			ON CONFLICT (wallet_id)
			DO UPDATE SET 
				balance = wallets.balance + EXCLUDED.balance,
				updated_at = NOW()
`
	_, err := executor.ExecContext(ctx,
		query,
		op.WalletID,
		op.Amount)
	if err != nil {
		return fmt.Errorf("failed update deposit:%w", err)
	}
	return nil
}

func (w WalletRepository) WithdrawBalance(ctx context.Context, op domain.WalletOperation) error {
	executor := w.executor.Get(ctx)

	query := `UPDATE wallets 
        	  SET balance = balance - $1 
              WHERE wallet_id = $2`

	_, err := executor.ExecContext(
		ctx,
		query,
		op.Amount,
		op.WalletID,
	)
	if err != nil {
		return fmt.Errorf("failed withdraw balance:%w", err)
	}
	return nil
}

func (w WalletRepository) GetWalletBalance(ctx context.Context, id uuid.UUID) (domain.Wallet, error) {
	executor := w.executor.Get(ctx)

	var wallet domain.Wallet

	query := `SELECT wallet_id, balance 
        	  FROM wallets 
       		  WHERE wallet_id = $1`

	err := executor.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(&wallet.ID,
		&wallet.Balance,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Wallet{}, domain.ErrWalletNotFound
		}
		return domain.Wallet{}, fmt.Errorf("failed get wallet: %w", err)

	}

	return wallet, nil
}
