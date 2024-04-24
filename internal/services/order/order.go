package order

import (
	apiError "accumulativeSystem/internal/errors/api"
	orderModel "accumulativeSystem/internal/models/order"
	orderRepository "accumulativeSystem/internal/repositories/order"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"net/http"
	"strconv"
	"time"
)

type OrderService interface {
	CreateOrder(order *orderModel.Order) (*orderModel.Order, error)
	GetOrderByOrderId(orderId int) (*orderModel.Order, error)
	GetOrdersByUserId(userId int) ([]*orderModel.Order, error)
	GetOrderByOrderIdAndUserID(orderId int, userId float64) (*orderModel.Order, error)
}

type orderService struct {
	repo orderRepository.OrderRepository
}

func NewOrderService(repo orderRepository.OrderRepository) OrderService {
	return &orderService{repo: repo}
}

func (s *orderService) CreateOrder(order *orderModel.Order) (*orderModel.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if !isValidLunaChecksum(order.OrderId) {
		return nil, apiError.NewError(http.StatusUnprocessableEntity, "not valid card number", nil)
	}

	existingOrder, err := s.GetOrderByOrderId(order.OrderId)

	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, apiError.NewError(http.StatusInternalServerError, "Internal Server Error", err)
		}
	}

	if existingOrder != nil && existingOrder.UserId != order.UserId {
		return nil, apiError.NewError(http.StatusConflict, "order already exists for another user", nil)
	}

	if existingOrder != nil {
		return nil, apiError.NewError(http.StatusOK, "order already exists", nil)
	}

	err = s.repo.CreateOrder(ctx, nil, order)

	//TODO добавить ошибку
	if err != nil {
		return nil, err
	}

	order, err = s.repo.GetOrderByOrderId(ctx, nil, order.OrderId)

	if err != nil {
		return nil, apiError.NewError(http.StatusInternalServerError, "Internal Server Error", err)
	}

	return order, nil
}

func (s *orderService) GetOrderByOrderId(orderId int) (*orderModel.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.repo.GetOrderByOrderId(ctx, nil, orderId)
}

func (s *orderService) GetOrdersByUserId(userId int) ([]*orderModel.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	orders, err := s.repo.GetOrdersByUserId(ctx, nil, userId)

	if err != nil {
		return nil, apiError.NewError(http.StatusInternalServerError, "Internal Server Error", err)
	}

	if len(orders) == 0 {
		return nil, apiError.NewError(http.StatusNoContent, "Internal Server Error", err)
	}

	return orders, nil
}

func (s *orderService) GetOrderByOrderIdAndUserID(orderId int, userId float64) (*orderModel.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.repo.GetOrderByOrderIdAndUserID(ctx, nil, orderId, userId)
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
