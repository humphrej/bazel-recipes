package testing

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4/pgxpool"
	"os"
)

//NewPool creates a pool from the DATABASE_URL environment variable - only for testing
func NewPool(ctx context.Context) (*pgxpool.Pool, error) {

	url, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		return nil, errors.New("need to set DATABASE_URL environment variable")
	}

	return pgxpool.Connect(ctx, url)
}
