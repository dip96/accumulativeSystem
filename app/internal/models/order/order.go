package order

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Order struct {
	Id               int              `json:"id"`
	UserId           int              `json:"user_id"`
	OrderId          int              `json:"order_id"`
	Accrual          float64          `json:"accrual"`
	WithdrawnBalance float64          `json:"withdrawn_balance"`
	Status           OrderStatus      `json:"status"`
	CreatedAt        pgtype.Timestamp `json:"created_at"`
}

type OrderStatus string

const (
	OrderStatusNew        OrderStatus = "NEW"
	OrderStatusRegistered OrderStatus = "REGISTERED"
	OrderStatusProcessing OrderStatus = "PROCESSING"
	OrderStatusInvalid    OrderStatus = "INVALID"
	OrderStatusProcessed  OrderStatus = "PROCESSED"
)

func GetOrderStatusByValue(value string) OrderStatus {
	switch value {
	case string(OrderStatusNew):
		return OrderStatusNew
	case string(OrderStatusRegistered):
		return OrderStatusRegistered
	case string(OrderStatusProcessing):
		return OrderStatusProcessing
	case string(OrderStatusInvalid):
		return OrderStatusInvalid
	case string(OrderStatusProcessed):
		return OrderStatusProcessed
	default:
		//TODO ERROR
		return ""
	}
}
