package storage

import "accumulativeSystem/internal/models/user"
import "accumulativeSystem/internal/models/order"

type StorageUserInterface interface {
	CreateUser(user user.User) error
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
