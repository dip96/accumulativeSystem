package postgres

import (
	balanceModel "accumulativeSystem/internal/models/balance"
	"context"
	"github.com/jackc/pgx/v5"
	"time"
)

func (s *Postgres) CreateBalance(balance *balanceModel.UserBalance) (*balanceModel.UserBalance, error) {
	//начинаем транзакцию
	tx, err := s.Pool.Begin(context.Background())
	if err != nil {
		return nil, err
	}

	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil {
			//TODO добавить логи
		}
	}(tx, context.Background())

	sqlQuery := "INSERT INTO user_balances (user_id, balance, withdrawn_balance) VALUES ($1,$2,$3)"
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//TODO не смог разобраться почему получаю ошибку в данном подходе
	//sqlQuery := "INSERT INTO users (login, password) VALUES (@login,@password)"
	//_, err := s.Pool.Exec(ctx, sqlQuery,
	//	sql.Named("login", user.Login),
	//	sql.Named("password", user.HashPassword),
	//)

	_, err = s.Pool.Exec(ctx, sqlQuery,
		balance.UserID,
		balance.Balance,
		balance.WithdrawnBalance,
	)

	if err != nil {
		return nil, err
	}

	userBalance, err := s.GetUserBalance(balance.UserID)

	if err != nil {
		return nil, err
	}

	return userBalance, nil
}

func (s *Postgres) GetUserBalance(userID int) (*balanceModel.UserBalance, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sqlSelect := "SELECT id, user_id, balance, withdrawn_balance FROM user_balances WHERE user_id = $1"
	var userBalance balanceModel.UserBalance
	err := s.Pool.QueryRow(ctx, sqlSelect, userID).Scan(&userBalance.ID, &userBalance.UserID, &userBalance.Balance, &userBalance.WithdrawnBalance)

	if err != nil {
		return nil, err
	}

	return &userBalance, nil
}
