package postgres

import (
	"accumulativeSystem/internal/errors/postgres"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"time"
)

type PoolWrapper struct {
	pool *pgxpool.Pool
}

func NewPoolWrapper(pool *pgxpool.Pool) *PoolWrapper {
	return &PoolWrapper{pool: pool}
}

func (pw *PoolWrapper) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	retryDelays := []time.Duration{1 * time.Second, 1 * time.Second, 1 * time.Second}
	for attempt, delay := range retryDelays {
		tag, err := pw.pool.Exec(ctx, sql, args...)

		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				if pgErr.Code == "23505" {
					return pgconn.CommandTag{}, postgres.New("duplicate login", pgErr)
				}
			}

			log.Printf("Err (attempt %d/%d): %v", attempt+1, len(retryDelays), err)
			time.Sleep(delay)
			continue
		}

		return tag, err
	}

	return pgconn.CommandTag{}, nil
}

func (pw *PoolWrapper) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return pw.pool.Query(ctx, sql, args...)
}

func (pw *PoolWrapper) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	row := pw.pool.QueryRow(ctx, sql, args...)
	return row
}

func (pw *PoolWrapper) Begin(ctx context.Context) (pgx.Tx, error) {
	return pw.pool.Begin(ctx)
}
