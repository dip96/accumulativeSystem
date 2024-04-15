package storage

import "accumulativeSystem/internal/model/user"
import "accumulativeSystem/internal/model/order"

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
