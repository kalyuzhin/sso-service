package app

import (
	grpcapp "github.com/kalyuzhin/sso-service/internal/app/grpc"
	"time"
)

// App – ...
type App struct {
	GRPCServer *grpcapp.GRPCApp
	TokenTTL   time.Duration
}

// New – ...
func New(port int) *App {
	grpcApp := grpcapp.NewGRPCApp(nil, port)

	return &App{
		GRPCServer: grpcApp,
	}
}
