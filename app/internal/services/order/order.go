package order

import (
	APIError "accumulativeSystem/internal/errors/api"
	"accumulativeSystem/internal/logger"
	orderModel "accumulativeSystem/internal/models/order"
	orderRepository "accumulativeSystem/internal/repositories/order"
	orderQueue "accumulativeSystem/internal/services/order/queue"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"net/http"
	"strconv"
	"time"
)

type OrderService interface {
	CreateOrder(order *orderModel.Order) (*orderModel.Order, error)
	SaveOrder(order *orderModel.Order) error
	GetOrderByOrderID(OrderID int) (*orderModel.Order, error)
	GetOrdersByUserID(UserID int) ([]*orderModel.Order, error)
	GetWithdrawalsByUserID(UserID int) ([]*orderModel.Order, error)
	GetOrderByOrderIDAndUserID(OrderID int, UserID float64) (*orderModel.Order, error)
}

type orderService struct {
	repo      orderRepository.OrderRepository
	chanOrder orderQueue.OrderQueueService
	logger    logger.Logger
}

func NewOrderService(repo orderRepository.OrderRepository, orderChan orderQueue.OrderQueueService, logger logger.Logger) OrderService {
	return &orderService{repo: repo, chanOrder: orderChan, logger: logger}
}

func (s *orderService) CreateOrder(order *orderModel.Order) (*orderModel.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if !isValidLunaChecksum(order.OrderID) {
		s.logger.Error("not valid card number")
		return nil, APIError.NewError(http.StatusUnprocessableEntity, "not valid card number", nil)
	}

	existingOrder, err := s.GetOrderByOrderID(order.OrderID)

	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			s.logger.Error(err.Error())
			return nil, APIError.NewError(http.StatusInternalServerError, "Internal Server Error", err)
		}
	}

	if existingOrder != nil && existingOrder.UserID != order.UserID {
		s.logger.Error("order already exists for another user")
		return nil, APIError.NewError(http.StatusConflict, "order already exists for another user", nil)
	}

	if existingOrder != nil {
		s.logger.Error("order already exists")
		return nil, APIError.NewError(http.StatusOK, "order already exists", nil)
	}

	err = s.repo.CreateOrder(ctx, nil, order)

	//TODO добавить ошибку
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	order, err = s.repo.GetOrderByOrderID(ctx, nil, order.OrderID)

	if err != nil {
		s.logger.Error(err.Error())
		return nil, APIError.NewError(http.StatusInternalServerError, "Internal Server Error", err)
	}

	s.chanOrder.EnqueueOrder(order.OrderID)

	return order, nil
}

func (s *orderService) GetOrderByOrderID(OrderID int) (*orderModel.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	order, err := s.repo.GetOrderByOrderID(ctx, nil, OrderID)

	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	return order, nil
}

func (s *orderService) GetOrdersByUserID(UserID int) ([]*orderModel.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	orders, err := s.repo.GetOrdersByUserID(ctx, nil, UserID)

	if err != nil {
		s.logger.Error(err.Error())
		return nil, APIError.NewError(http.StatusInternalServerError, "Internal Server Error", err)
	}

	if len(orders) == 0 {
		s.logger.Error("Not found orders")
		return nil, APIError.NewError(http.StatusNoContent, "Not found orders", err)
	}

	return orders, nil
}

func (s *orderService) GetWithdrawalsByUserID(UserID int) ([]*orderModel.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	orders, err := s.repo.GetWithdrawalsByUserID(ctx, nil, UserID)

	if err != nil {
		s.logger.Error(err.Error())
		return nil, APIError.NewError(http.StatusInternalServerError, "Internal Server Error", err)
	}

	if len(orders) == 0 {
		s.logger.Error("Not found orders")
		return nil, APIError.NewError(http.StatusNoContent, "Not found orders", err)
	}

	return orders, nil
}

func (s *orderService) GetOrderByOrderIDAndUserID(OrderID int, UserID float64) (*orderModel.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	order, err := s.repo.GetOrderByOrderIDAndUserID(ctx, nil, OrderID, UserID)

	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	return order, nil
}

func (s *orderService) SaveOrder(order *orderModel.Order) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := s.repo.Save(ctx, nil, order)

	if err != nil {
		s.logger.Error(err.Error())
		return err
	}

	return nil
}

// TODO перенести в lib???
func isValidLunaChecksum(creditCardNumber int) bool {
	var sum int
	var isEven = false

	strCardNumber := strconv.Itoa(creditCardNumber)
	for i := len(strCardNumber) - 1; i >= 0; i-- {
		digit, _ := strconv.Atoi(string(strCardNumber[i]))
		if isEven {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
		isEven = !isEven
	}

	return sum%10 == 0
}
