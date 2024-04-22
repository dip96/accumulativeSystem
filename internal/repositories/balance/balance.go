package balance

import (
	balanceModel "accumulativeSystem/internal/models/balance"
	storage "accumulativeSystem/internal/storage/postgres"
	"context"
	"github.com/jackc/pgx/v5"
)

type BalanceRepository interface {
	CreateBalance(ctx context.Context, balance *balanceModel.UserBalance) (*balanceModel.UserBalance, error)
	GetUserBalance(ctx context.Context, userID int) (*balanceModel.UserBalance, error)
	UpdateUserBalance(ctx context.Context, balance *balanceModel.UserBalance) error
	Begin(ctx context.Context) (pgx.Tx, error)
}

type balanceRepository struct {
	db storage.Storage
}

func NewBalanceRepository(storage storage.Storage) BalanceRepository {
	return &balanceRepository{db: storage}
}

func (r *balanceRepository) CreateBalance(ctx context.Context, balance *balanceModel.UserBalance) (*balanceModel.UserBalance, error) {
	sqlQuery := "INSERT INTO user_balances (user_id, balance, withdrawn_balance) VALUES ($1,$2,$3)"

	_, err := r.db.Exec(ctx, sqlQuery,
		balance.UserID,
		balance.Balance,
		balance.WithdrawnBalance,
	)

	if err != nil {
		return nil, err
	}

	userBalance, err := r.GetUserBalance(ctx, balance.UserID)

	if err != nil {
		return nil, err
	}

	return userBalance, nil
}

func (r *balanceRepository) GetUserBalance(ctx context.Context, userID int) (*balanceModel.UserBalance, error) {
	sqlSelect := "SELECT id, user_id, balance, withdrawn_balance FROM user_balances WHERE user_id = $1"
	var userBalance balanceModel.UserBalance
	err := r.db.QueryRow(ctx, sqlSelect, userID).Scan(&userBalance.ID, &userBalance.UserID, &userBalance.Balance, &userBalance.WithdrawnBalance)
	if err != nil {
		return nil, err
	}

	return &userBalance, nil
}

func (r *balanceRepository) UpdateUserBalance(ctx context.Context, balance *balanceModel.UserBalance) error {
	sqlUpdate := "UPDATE user_balances SET balance = $1, withdrawn_balance = $2 WHERE user_id = $3"

	_, err := r.db.Exec(ctx, sqlUpdate, balance.Balance, balance.WithdrawnBalance, balance.UserID)
	if err != nil {
		return err
	}

	return nil
}

func (r *balanceRepository) Begin(ctx context.Context) (pgx.Tx, error) {
	return r.db.Begin(ctx)
}
