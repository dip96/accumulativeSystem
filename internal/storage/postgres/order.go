package postgres

import (
	orderModel "accumulativeSystem/internal/models/order"
	"context"
	//"github.com/jackc/pgx/v5"
	"time"
)

func (s *Postgres) CreateOrder(userId float64, orderId int) (*orderModel.Order, error) {
	//начинаем транзакцию
	//tx, err := s.Pool.Begin(context.Background())
	//if err != nil {
	//	return nil, err
	//}
	//
	//defer func(tx pgx.Tx, ctx context.Context) {
	//	err := tx.Rollback(ctx)
	//	if err != nil {
	//		//TODO добавить логи
	//	}
	//}(tx, context.Background())

	sqlQuery := "INSERT INTO orders (user_id, order_id, status, accrual) VALUES ($1,$2, $3, $4)"
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := s.Pool.Exec(ctx, sqlQuery,
		userId,
		orderId,
		orderModel.OrderStatusNew,
		0,
	)

	if err != nil {
		return nil, err
	}

	order, err := s.GetOrderByOrderId(orderId)

	if err != nil {
		return nil, err
	}

	return order, nil
}

func (s *Postgres) GetOrderByOrderId(orderId int) (*orderModel.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sqlSelect := "SELECT id, user_id, order_id, accrual, status, created_at FROM orders WHERE order_id = $1"
	var order orderModel.Order
	err := s.Pool.QueryRow(ctx, sqlSelect, orderId).Scan(&order.Id, &order.UserId, &order.OrderId, &order.Accrual, &order.Status, &order.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (s *Postgres) GetOrderByOrderIdAndUserID(orderId int, userId float64) (*orderModel.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sqlSelect := "SELECT id, user_id, order_id, accrual, status, created_at FROM orders WHERE order_id = $1 AND user_id = $2"
	var order orderModel.Order
	err := s.Pool.QueryRow(ctx, sqlSelect, orderId, userId).Scan(&order.Id, &order.UserId, &order.OrderId, &order.Accrual, &order.Status, &order.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &order, nil
}
