package order

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Order struct {
	Id        int              `json:"id"`
	UserId    int              `json:"user_id"`
	OrderId   int              `json:"order_id"`
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
