package balance

import (
	balanceModel "accumulativeSystem/internal/models/balance"
	"accumulativeSystem/internal/storage"
	"context"
	"github.com/jackc/pgx/v5"
)

type BalanceRepository interface {
	CreateBalance(ctx context.Context, tx pgx.Tx, balance *balanceModel.UserBalance) error
	GetUserBalance(ctx context.Context, tx pgx.Tx, userID int) (*balanceModel.UserBalance, error)
	UpdateUserBalance(ctx context.Context, tx pgx.Tx, balance *balanceModel.UserBalance) error
	Begin(ctx context.Context) (pgx.Tx, error)
}

type balanceRepository struct {
	db storage.Storage
}

func NewBalanceRepository(storage storage.Storage) BalanceRepository {
	return &balanceRepository{db: storage}
}

func (r *balanceRepository) CreateBalance(ctx context.Context, tx pgx.Tx, balance *balanceModel.UserBalance) error {
	sqlQuery := "INSERT INTO user_balances (user_id, balance, withdrawn_balance) VALUES ($1,$2,$3)"

	if tx == nil {
		_, err := r.db.Exec(ctx, sqlQuery, balance.UserID, balance.Balance, balance.WithdrawnBalance)
		return err
	}

	_, err := tx.Exec(ctx, sqlQuery, balance.UserID, balance.Balance, balance.WithdrawnBalance)
	return err
}

func (r *balanceRepository) GetUserBalance(ctx context.Context, tx pgx.Tx, UserID int) (*balanceModel.UserBalance, error) {
	sqlSelect := "SELECT id, user_id, balance, withdrawn_balance FROM user_balances WHERE user_id = $1"
	var userBalance balanceModel.UserBalance

	if tx == nil {
		err := r.db.QueryRow(ctx, sqlSelect, UserID).Scan(&userBalance.ID, &userBalance.UserID, &userBalance.Balance, &userBalance.WithdrawnBalance)
		if err != nil {
			return nil, err
		}
		return &userBalance, nil
	}

	err := tx.QueryRow(ctx, sqlSelect, UserID).Scan(&userBalance.ID, &userBalance.UserID, &userBalance.Balance, &userBalance.WithdrawnBalance)
	if err != nil {
		return nil, err
	}

	return &userBalance, nil
}

func (r *balanceRepository) UpdateUserBalance(ctx context.Context, tx pgx.Tx, balance *balanceModel.UserBalance) error {
	sqlUpdate := "UPDATE user_balances SET balance = $1, withdrawn_balance = $2 WHERE user_id = $3"

	if tx == nil {
		_, err := r.db.Exec(ctx, sqlUpdate, balance.Balance, balance.WithdrawnBalance, balance.UserID)
		return err
	}

	_, err := tx.Exec(ctx, sqlUpdate, balance.Balance, balance.WithdrawnBalance, balance.UserID)
	return err
}

func (r *balanceRepository) Begin(ctx context.Context) (pgx.Tx, error) {
	return r.db.Begin(ctx)
}
