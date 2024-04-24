package balance

import (
	apiError "accumulativeSystem/internal/errors/api"
	balanceModel "accumulativeSystem/internal/models/balance"
	orderModel "accumulativeSystem/internal/models/order"
	balanceRepository "accumulativeSystem/internal/repositories/balance"
	"accumulativeSystem/internal/repositories/order"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"
	"net/http"
	"time"
)

type BalanceService interface {
	CreateBalance(balance *balanceModel.UserBalance) (*balanceModel.UserBalance, error)
	GetUserBalance(userID int) (*balanceModel.UserBalance, error)
	UpdateUserBalance(balance *balanceModel.UserBalance) error
	WithdrawBalance(userID int, orderID int, sum float64) error
}

type balanceService struct {
	repo      balanceRepository.BalanceRepository
	repoOrder order.OrderRepository
}

func NewBalanceService(repo balanceRepository.BalanceRepository, repoOrder order.OrderRepository) BalanceService {
	return &balanceService{repo: repo, repoOrder: repoOrder}
}

func (s *balanceService) CreateBalance(balance *balanceModel.UserBalance) (*balanceModel.UserBalance, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := s.repo.CreateBalance(ctx, nil, balance)

	usBalance, err := s.repo.GetUserBalance(ctx, nil, balance.UserID)

	if err != nil {
		return nil, err
	}

	//TODO добавить ошибку
	if err != nil {
		return nil, err
	}

	return usBalance, nil
}

func (s *balanceService) GetUserBalance(userID int) (*balanceModel.UserBalance, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	balance, err := s.repo.GetUserBalance(ctx, nil, userID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apiError.NewError(http.StatusInternalServerError, "no balance information", err)
		}
	}

	return balance, nil
}

func (s *balanceService) UpdateUserBalance(balance *balanceModel.UserBalance) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.repo.UpdateUserBalance(ctx, nil, balance)
}

func (s *balanceService) WithdrawBalance(userID int, orderID int, sum float64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sumDecimal := decimal.NewFromFloat(sum)

	// Начинаем транзакцию
	tx, err := s.repo.Begin(ctx)
	if err != nil {
		return apiError.NewError(http.StatusInternalServerError, "Internal Server Error", err)
	}

	defer func() {
		//TODO интересный момент в случаи паники, err == nil
		if err != nil {
			tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
			if err != nil {
				return
			}
		}
	}()

	userBalance, err := s.repo.GetUserBalance(ctx, tx, userID)
	userBalanceDecimal := decimal.NewFromFloat(userBalance.Balance)

	if err != nil {
		return apiError.NewError(http.StatusInternalServerError, "Internal Server Error", nil)
	}

	//Проверяем, достаточно ли средств
	if userBalanceDecimal.LessThan(sumDecimal) {
		return apiError.NewError(http.StatusPaymentRequired, "insufficient funds", nil)
	}

	// Создаем новый заказ
	var order orderModel.Order
	order.OrderId = orderID
	order.UserId = userID
	err = s.repoOrder.CreateOrder(ctx, tx, &order)
	if err != nil {
		return apiError.NewError(http.StatusInternalServerError, "Internal Server Error", nil)
	}

	// Обновляем баланс пользователя
	newBalance := userBalanceDecimal.Sub(sumDecimal)
	newWithdrawnBalance := decimal.NewFromFloat(userBalance.WithdrawnBalance).Add(sumDecimal)

	userBalance.Balance, _ = newBalance.Float64()
	userBalance.WithdrawnBalance, _ = newWithdrawnBalance.Float64()
	err = s.repo.UpdateUserBalance(ctx, tx, userBalance)
	if err != nil {
		return apiError.NewError(http.StatusInternalServerError, "Internal Server Error", nil)
	}

	return nil
}
