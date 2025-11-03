package grpcapp

import (
	"fmt"
	authgrpc "github.com/kalyuzhin/sso-service/internal/grpc/auth"
	"google.golang.org/grpc"
	"net"
)

// GRPCApp – ...
type GRPCApp struct {
	grpcServer *grpc.Server
	port       int
}

// NewGRPCApp – ...
func NewGRPCApp(log any, port int) *GRPCApp {
	server := grpc.NewServer()
	authgrpc.Register(server)

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
