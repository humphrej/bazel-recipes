package grpc

import (
	"context"
	"github.com/golang/glog"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	"turps/api"
	"turps/internal"
	"turps/internal/postgres"
)

func NewTurpsServer(ctx context.Context, dbUrl string) *grpc.Server {
	pool, err := pgxpool.Connect(ctx, dbUrl)
	if err != nil {
		glog.Fatalf("failed to create connection pool: %v", err)
	}
	var storage internal.ChangeListRepository
	storage = postgres.ChangeListStorage{Pool: pool}

	service := NewServer(storage)
	s := grpc.NewServer()
	api.RegisterTurpsServer(s, service)
	return s
}
