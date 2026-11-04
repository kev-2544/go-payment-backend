package repository

import (
    "context"
    "fmt"

    "github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresPool(
    ctx context.Context,
) (*pgxpool.Pool, error) {

    connString := "postgres://postgres:Manoti123456@localhost:5432/payments_db"

    pool, err := pgxpool.New(ctx, connString)
    if err != nil {
        return nil, fmt.Errorf("connect to postgres: %w", err)
    }

    if err := pool.Ping(ctx); err != nil {
        pool.Close()
        return nil, fmt.Errorf("ping postgres: %w", err)
    }

    return pool, nil
}