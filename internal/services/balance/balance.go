package balance

import (
	balanceModel "accumulativeSystem/internal/models/balance"
	orderModel "accumulativeSystem/internal/models/order"
	balanceRepository "accumulativeSystem/internal/repositories/balance"
	"accumulativeSystem/internal/repositories/order"
	"context"
	"errors"
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

func NewBalanceService(repo balanceRepository.BalanceRepository) BalanceService {
	return &balanceService{repo: repo}
}

func (s *balanceService) CreateBalance(balance *balanceModel.UserBalance) (*balanceModel.UserBalance, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	usBalance, err := s.repo.CreateBalance(ctx, balance)

	//TODO добавить ошибку
	if err != nil {
		return nil, err
	}

	return usBalance, nil
}

func (s *balanceService) GetUserBalance(userID int) (*balanceModel.UserBalance, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.repo.GetUserBalance(ctx, userID)
}

func (s *balanceService) UpdateUserBalance(balance *balanceModel.UserBalance) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.repo.UpdateUserBalance(ctx, balance)
}

func (s *balanceService) WithdrawBalance(userID int, orderID int, sum float64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//TODO ТРАНЗАКЦИЯ

	userBalance, err := s.repo.GetUserBalance(ctx, userID)

	if err != nil {
		return err
	}

	//Проверяем, достаточно ли средств
	if userBalance.Balance < sum {
		return errors.New("insufficient funds") //TODO добавить кастомную ошибку
	}

	// Создаем новый заказ
	var order orderModel.Order
	order.OrderId = orderID
	order.UserId = userID
	err = s.repoOrder.CreateOrder(ctx, &order)
	if err != nil {
		return err
	}

	// Обновляем баланс пользователя
	userBalance.Balance -= sum
	userBalance.WithdrawnBalance += sum
	err = s.repo.UpdateUserBalance(ctx, userBalance)
	if err != nil {
		return err
	}

	return nil
}
