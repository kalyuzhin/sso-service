package auth

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"regexp"

	errorpkg "github.com/kalyuzhin/sso-service/internal/error"
	ssov1 "github.com/kalyuzhin/sso-service/internal/pkg/pb"
)

const (
	emptyAppIDValue = 0
)

type serverAPI struct {
	ssov1.UnimplementedAuthServer
}

// Register – ...
func Register(gRPC *grpc.Server) {
	ssov1.RegisterAuthServer(gRPC, &serverAPI{})
}

type implementation interface {
}

// Register – ...
func (s *serverAPI) Register(ctx context.Context, req *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {
	err := validateEmailAndPassword(req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &ssov1.RegisterResponse{}, nil
}

// Login – ...
func (s *serverAPI) Login(ctx context.Context, req *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {
	email := req.GetEmail()
	password := req.GetPassword()

	err := validateEmailAndPassword(email, password)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	appID := req.GetAddId()
	if appID == emptyAppIDValue {
		return nil, status.Error(codes.InvalidArgument, "app id is required")
	}

	return &ssov1.LoginResponse{}, nil
}

func validateEmailAndPassword(email, password string) error {
	if email == "" {
		return errorpkg.New("email is required")
	}

	if password == "" {
		return errorpkg.New("password is required")
	}

	if len(password) < 7 {
		return errorpkg.New("password length should be 8 characters at least")
	}

	re := regexp.MustCompile(`w+@w+\.w+`)
	if ok := re.MatchString(email); !ok {
		return errorpkg.New("email doesn't match")
	}

	return nil
}
