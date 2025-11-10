package auth

import (
	"context"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"

	"github.com/kalyuzhin/sso-service/internal/model"
)

func getAuxiliaryParams(ctx context.Context) (userAgent []string, ip string, err error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return userAgent, ip, model.ErrMetaData
	}
	userAgent = md.Get(model.UserAgent)
	if len(userAgent) < 1 {
		return userAgent, ip, model.ErrUserAgent
	}

	p, ok := peer.FromContext(ctx)
	if !ok {
		return userAgent, ip, model.ErrIP
	}
	ip = p.Addr.String()

	return userAgent, ip, nil
}
