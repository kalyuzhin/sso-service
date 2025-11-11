package auth

import (
	"context"
	errorpkg "github.com/kalyuzhin/sso-service/internal/error"
	"github.com/kalyuzhin/sso-service/internal/model"
	ssov1 "github.com/kalyuzhin/sso-service/internal/pkg/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/mail"
)

const (
	emptyAppIDValue = 0
)

type serverAPI struct {
	ssov1.UnimplementedAuthServer
	service implementation
}

// Register – ...
func Register(gRPC *grpc.Server, service implementation) {
	ssov1.RegisterAuthServer(gRPC, &serverAPI{service: service})
}

type implementation interface {
	Register(ctx context.Context, email, pswd string) (userID int64, err error)
	Login(ctx context.Context, email, pswd string, appID int32, params model.UserRequestParams) (token string, err error)
}

// Register – ...
func (s *serverAPI) Register(ctx context.Context, req *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {
	email := req.GetEmail()
	password := req.GetPassword()

	err := validateEmailAndPassword(email, password)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	userID, err := s.service.Register(ctx, email, password)
	if err != nil || userID == 0 {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.RegisterResponse{UserId: userID}, nil
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

	userAgent, ip, err := getAuxiliaryParams(ctx)
	if err != nil {
		return nil, status.Error(codes.DataLoss, "invalid request")
	}

	token, err := s.service.Login(ctx, email, password, appID, model.UserRequestParams{
		IP:        ip,
		UserAgent: concatString(userAgent),
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &ssov1.LoginResponse{Token: token}, nil
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

	a, err := mail.ParseAddress(email)
	if err != nil || a.Address != email {
		return errorpkg.New("email validation failed")
	}

	return nil
}
