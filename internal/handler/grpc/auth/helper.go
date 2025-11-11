package auth

import (
	"context"
	"strings"

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

func concatString(input []string) string {
	sb := strings.Builder{}
	for _, s := range input {
		sb.WriteString(s)
	}

	return sb.String()
}
