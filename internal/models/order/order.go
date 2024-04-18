package order

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Order struct {
	Id        int64            `json:"id"`
	UserId    float64          `json:"user_id"`
	OrderId   int64            `json:"order_id"`
	Accrual   float64          `json:"accrual"`
	Status    OrderStatus      `json:"status"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
}

type OrderStatus string

const (
	OrderStatusNew        OrderStatus = "NEW"
	OrderStatusRegistered OrderStatus = "REGISTERED"
	OrderStatusProcessing OrderStatus = "PROCESSING"
	OrderStatusInvalid    OrderStatus = "INVALID"
	OrderStatusProcessed  OrderStatus = "PROCESSED"
)
