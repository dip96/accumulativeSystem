package order

import (
	orderModel "accumulativeSystem/internal/models/order"
	orderRepository "accumulativeSystem/internal/repositories/order"
	"context"
	"errors"
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
		return nil, errors.New("not valid card number")
	}

	existingOrder, err := s.GetOrderByOrderId(order.OrderId)
	if existingOrder != nil && existingOrder.UserId != order.UserId {
		return nil, errors.New("order already exists for another user")
	}

	if err == nil && existingOrder != nil {
		return nil, errors.New("order already exists")
	}

	err = s.repo.CreateOrder(ctx, order)

	//TODO добавить ошибку
	if err != nil {
		return nil, err
	}

	order, err = s.repo.GetOrderByOrderId(ctx, order.OrderId)

	if err != nil {
		return nil, err
	}

	return order, nil
}

func (s *orderService) GetOrderByOrderId(orderId int) (*orderModel.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.repo.GetOrderByOrderId(ctx, orderId)
}

func (s *orderService) GetOrdersByUserId(userId int) ([]*orderModel.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.repo.GetOrdersByUserId(ctx, userId)
}

func (s *orderService) GetOrderByOrderIdAndUserID(orderId int, userId float64) (*orderModel.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.repo.GetOrderByOrderIdAndUserID(ctx, orderId, userId)
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
