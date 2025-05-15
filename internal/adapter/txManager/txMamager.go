package txManager

import (
	"context"
	"database/sql"
	"fmt"
)

type TxKey string

const (
	TxKeyName TxKey = "tx"
)

type TxManager struct {
	db *sql.DB
}

func NewTxManager(db *sql.DB) *TxManager {
	return &TxManager{
		db: db,
	}
}

func (t TxManager) Do(ctx context.Context, fn func(ctx context.Context) error) (err error) {
	tx, err := t.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() {
		if err != nil {
			if e := tx.Rollback(); e != nil {
				err = fmt.Errorf("rollback tx: %w:%w", err, e)

				return
			}

			return
		}

		if e := tx.Commit(); e != nil {
			err = fmt.Errorf("commit tx: %w:%w", err, e)
		}
	}()

	ctxWithTx := context.WithValue(ctx, TxKeyName, tx)

	if err := fn(ctxWithTx); err != nil {
		return err
	}

	return nil
}
