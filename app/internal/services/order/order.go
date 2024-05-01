package order

import (
	apiError "accumulativeSystem/internal/errors/api"
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
	GetOrderByOrderId(orderId int) (*orderModel.Order, error)
	GetOrdersByUserId(userId int) ([]*orderModel.Order, error)
	GetWithdrawalsByUserId(userId int) ([]*orderModel.Order, error)
	GetOrderByOrderIdAndUserID(orderId int, userId float64) (*orderModel.Order, error)
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

	if !isValidLunaChecksum(order.OrderId) {
		s.logger.Error("not valid card number")
		return nil, apiError.NewError(http.StatusUnprocessableEntity, "not valid card number", nil)
	}

	existingOrder, err := s.GetOrderByOrderId(order.OrderId)

	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			s.logger.Error(err.Error())
			return nil, apiError.NewError(http.StatusInternalServerError, "Internal Server Error", err)
		}
	}

	if existingOrder != nil && existingOrder.UserId != order.UserId {
		s.logger.Error("order already exists for another user")
		return nil, apiError.NewError(http.StatusConflict, "order already exists for another user", nil)
	}

	if existingOrder != nil {
		s.logger.Error("order already exists")
		return nil, apiError.NewError(http.StatusOK, "order already exists", nil)
	}

	err = s.repo.CreateOrder(ctx, nil, order)

	//TODO добавить ошибку
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	order, err = s.repo.GetOrderByOrderId(ctx, nil, order.OrderId)

	if err != nil {
		s.logger.Error(err.Error())
		return nil, apiError.NewError(http.StatusInternalServerError, "Internal Server Error", err)
	}

	s.chanOrder.EnqueueOrder(order.OrderId)

	return order, nil
}

func (s *orderService) GetOrderByOrderId(orderId int) (*orderModel.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	order, err := s.repo.GetOrderByOrderId(ctx, nil, orderId)

	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	return order, nil
}

func (s *orderService) GetOrdersByUserId(userId int) ([]*orderModel.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	orders, err := s.repo.GetOrdersByUserId(ctx, nil, userId)

	if err != nil {
		s.logger.Error(err.Error())
		return nil, apiError.NewError(http.StatusInternalServerError, "Internal Server Error", err)
	}

	if len(orders) == 0 {
		s.logger.Error("Not found orders")
		return nil, apiError.NewError(http.StatusNoContent, "Not found orders", err)
	}

	return orders, nil
}

func (s *orderService) GetWithdrawalsByUserId(userId int) ([]*orderModel.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	orders, err := s.repo.GetWithdrawalsByUserId(ctx, nil, userId)

	if err != nil {
		s.logger.Error(err.Error())
		return nil, apiError.NewError(http.StatusInternalServerError, "Internal Server Error", err)
	}

	if len(orders) == 0 {
		s.logger.Error("Not found orders")
		return nil, apiError.NewError(http.StatusNoContent, "Not found orders", err)
	}

	return orders, nil
}

func (s *orderService) GetOrderByOrderIdAndUserID(orderId int, userId float64) (*orderModel.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	order, err := s.repo.GetOrderByOrderIdAndUserID(ctx, nil, orderId, userId)

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
