package auth

import (
	"context"

	"google.golang.org/grpc"

	ssov1 "github.com/kalyuzhin/sso-service/internal/pkg/pb"
)

type serverAPI struct {
	ssov1.UnimplementedAuthServer
}

// Register – ...
func Register(gRPC *grpc.Server) {
	ssov1.RegisterAuthServer(gRPC, &serverAPI{})
}

// Register – ...
func (s *serverAPI) Register(ctx context.Context, req *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {
	panic("implement me")
}

// Login – ...
func (s *serverAPI) Login(ctx context.Context, req *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {
	panic("implement me")
}
