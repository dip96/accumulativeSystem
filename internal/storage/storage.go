package storage

import (
	"accumulativeSystem/internal/models/order"
	userModel "accumulativeSystem/internal/models/user"
)

type StorageUserInterface interface {
	CreateUser(login string, password []byte) (*userModel.User, error)
	GetUser(login string) (*userModel.User, error)
	GetUserPassword(login string) (*userModel.User, error)
}

type StorageOrderInterface interface {
	CreateOrder(user order.Order) error
	UpdateOrder(order.Order) error
	GetOrderById(id int) (order.Order, error)
	GetOrderByOrderId(orderId int) (order.Order, error)
	GetOrderByUserId(userId int) (order.Order, error)
}

type StorageInterface interface {
	StorageUserInterface
	StorageOrderInterface
}
