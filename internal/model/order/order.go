package order

import "github.com/jackc/pgx/v5/pgtype"

type Order struct {
	Id        int
	UserId    int
	OrderId   int
	Accrual   pgtype.Numeric
	Status    pgtype.EnumCodec
	CreatedAt pgtype.Timestamp
}
