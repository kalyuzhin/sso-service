package grpcapp

import (
	"fmt"
	"net"

	"google.golang.org/grpc"

	authgrpc "github.com/kalyuzhin/sso-service/internal/handler/grpc/auth"
	"github.com/kalyuzhin/sso-service/internal/service/auth"
)

// GRPCApp – ...
type GRPCApp struct {
	grpcServer *grpc.Server
	port       int
}

// NewGRPCApp – ...
func NewGRPCApp(authService *auth.Auth, port int) *GRPCApp {
	server := grpc.NewServer()
	authgrpc.Register(server, authService)

	return &GRPCApp{
		grpcServer: server,
		port:       port,
	}
}

// Run – ...
func (a *GRPCApp) Run() error {
	const op = "grpcapp.Run"

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return err
	}

	if err = a.grpcServer.Serve(l); err != nil {
		return err
	}

	return nil
}

// Stop – ...
func (a *GRPCApp) Stop() {
	a.grpcServer.GracefulStop()
}
