package auth

import (
	"context"
	"github.com/kalyuzhin/sso-service/internal/model"
	ssov1 "github.com/kalyuzhin/sso-service/internal/pkg/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	GetPublicKeys(_ context.Context) *model.JWKS
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

func (s *serverAPI) GetPublicKeys(ctx context.Context, req *ssov1.GetPublicKeyRequest) (*ssov1.GetPublicKeyResponse, error) {
	jwks := s.service.GetPublicKeys(ctx)
	resp := ssov1.GetPublicKeyResponse{Keys: make([]*ssov1.Jwk, 0, len(jwks.Keys))}

	for _, jwk := range jwks.Keys {
		resp.Keys = append(resp.Keys, &ssov1.Jwk{
			Kty: jwk.Kty,
			Alg: jwk.Alg,
			N:   jwk.N,
			E:   jwk.E,
			Use: jwk.Use,
			Kid: jwk.KID,
		})
	}

	return &resp, nil
}
