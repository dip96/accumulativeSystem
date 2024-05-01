package postgres

import (
	"accumulativeSystem/internal/config"
	postgresError "accumulativeSystem/internal/errors/postgres"
	storageInterface "accumulativeSystem/internal/storage"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"time"
)

type storage struct {
	pool *pgxpool.Pool
}

func NewDb(cnf config.ConfigInstance) (storageInterface.Storage, error) {
	pool, err := pgxpool.New(context.Background(), cnf.GetDatabaseURI())
	if err != nil {
		return nil, postgresError.New("error while initializing postgres pool", err)
	}
	return &storage{pool: pool}, nil
}

func (s *storage) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	retryDelays := []time.Duration{1 * time.Second, 1 * time.Second, 1 * time.Second}
	for attempt, delay := range retryDelays {
		tag, err := s.pool.Exec(ctx, sql, args...)

		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				if pgErr.Code == "23505" {
					return pgconn.CommandTag{}, postgresError.New("duplicate login", pgErr)
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

func (s *storage) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return s.pool.Query(ctx, sql, args...)
}

func (s *storage) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	row := s.pool.QueryRow(ctx, sql, args...)
	return row
}

func (s *storage) Begin(ctx context.Context) (pgx.Tx, error) {
	return s.pool.Begin(ctx)
}

func (s *storage) Acquire(ctx context.Context) (*pgxpool.Conn, error) {
	return s.pool.Acquire(ctx)
}

func (s *storage) Close() {
	s.pool.Close()
}
