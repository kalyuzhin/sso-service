package app

import (
	"context"
	"time"

	grpcapp "github.com/kalyuzhin/sso-service/internal/app/grpc"
	"github.com/kalyuzhin/sso-service/internal/config"
	"github.com/kalyuzhin/sso-service/internal/service/auth"
	"github.com/kalyuzhin/sso-service/internal/storage/postgresql"
)

// App – ...
type App struct {
	GRPCServer *grpcapp.GRPCApp
	TokenTTL   time.Duration
}

// New – ...
func New(ctx context.Context, dsn string, port int, cfg config.Config) (*App, error) {
	storage, err := postgresql.NewDB(ctx, dsn)
	if err != nil {
		return nil, err
	}

	authService := auth.New(storage, storage, storage, storage, cfg)
	grpcApp := grpcapp.NewGRPCApp(authService, port)

	return &App{
		GRPCServer: grpcApp,
	}, nil
}
