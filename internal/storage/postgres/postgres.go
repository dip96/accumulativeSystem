package postgres

import (
	"accumulativeSystem/internal/config"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	Pool *PoolWrapper
}

func NewDb() *Postgres {
	cnf := config.MustLoad()

	pool, err := pgxpool.New(context.Background(), cnf.DatabaseUri)
	if err != nil {
		panic(err)
	}
	wrappedPool := NewPoolWrapper(pool)
	return &Postgres{Pool: wrappedPool}
}
