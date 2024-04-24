package order

import (
	orderModel "accumulativeSystem/internal/models/order"
	"accumulativeSystem/internal/storage"
	"context"
	"github.com/jackc/pgx/v5"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, tx pgx.Tx, order *orderModel.Order) error
	GetOrderByOrderId(ctx context.Context, tx pgx.Tx, orderId int) (*orderModel.Order, error)
	GetOrdersByUserId(ctx context.Context, tx pgx.Tx, userId int) ([]*orderModel.Order, error)
	GetOrderByOrderIdAndUserID(ctx context.Context, tx pgx.Tx, orderId int, userId float64) (*orderModel.Order, error)
	Begin(ctx context.Context) (pgx.Tx, error)
}

type orderRepository struct {
	db storage.Storage
}

func NewOrderRepository(storage storage.Storage) OrderRepository {
	return &orderRepository{db: storage}
}

func (o *orderRepository) CreateOrder(ctx context.Context, tx pgx.Tx, order *orderModel.Order) error {
	sqlQuery := "INSERT INTO orders (user_id, order_id, status, accrual) VALUES ($1,$2, $3, $4)"

	if tx == nil {
		_, err := o.db.Exec(ctx, sqlQuery, order.UserId, order.OrderId, orderModel.OrderStatusNew, 0)
		return err
	}

	_, err := tx.Exec(ctx, sqlQuery, order.UserId, order.OrderId, orderModel.OrderStatusNew, 0)
	return err
}
func (o *orderRepository) GetOrderByOrderId(ctx context.Context, tx pgx.Tx, orderId int) (*orderModel.Order, error) {
	sqlSelect := "SELECT id, user_id, order_id, accrual, status, created_at FROM orders WHERE order_id = $1"
	var order orderModel.Order

	if tx == nil {
		err := o.db.QueryRow(ctx, sqlSelect, orderId).Scan(&order.Id, &order.UserId, &order.OrderId, &order.Accrual, &order.Status, &order.CreatedAt)
		if err != nil {
			return nil, err
		}
		return &order, nil
	}

	err := tx.QueryRow(ctx, sqlSelect, orderId).Scan(&order.Id, &order.UserId, &order.OrderId, &order.Accrual, &order.Status, &order.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (o *orderRepository) GetOrdersByUserId(ctx context.Context, tx pgx.Tx, userId int) ([]*orderModel.Order, error) {
	sqlSelect := "SELECT id, user_id, order_id, accrual, status, created_at FROM orders WHERE user_id = $1 ORDER BY created_at ASC"

	rows, err := o.db.Query(ctx, sqlSelect, userId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*orderModel.Order
	for rows.Next() {
		order := &orderModel.Order{}
		err = rows.Scan(&order.Id, &order.UserId, &order.OrderId, &order.Accrual, &order.Status, &order.CreatedAt)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (o *orderRepository) GetOrderByOrderIdAndUserID(ctx context.Context, tx pgx.Tx, orderId int, userId float64) (*orderModel.Order, error) {
	sqlSelect := "SELECT id, user_id, order_id, accrual, status, created_at FROM orders WHERE order_id = $1 AND user_id = $2"
	var order orderModel.Order

	if tx == nil {
		err := o.db.QueryRow(ctx, sqlSelect, orderId, userId).Scan(&order.Id, &order.UserId, &order.OrderId, &order.Accrual, &order.Status, &order.CreatedAt)
		if err != nil {
			return nil, err
		}
		return &order, nil
	}

	err := tx.QueryRow(ctx, sqlSelect, orderId, userId).Scan(&order.Id, &order.UserId, &order.OrderId, &order.Accrual, &order.Status, &order.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (o *orderRepository) Begin(ctx context.Context) (pgx.Tx, error) {
	return o.db.Begin(ctx)
}
