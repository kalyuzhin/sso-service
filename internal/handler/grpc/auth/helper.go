package auth

import (
	"context"
	errorpkg "github.com/kalyuzhin/sso-service/internal/error"
	"net/mail"
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

func validateEmailAndPassword(email, password string) error {
	if email == "" {
		return errorpkg.New("email is required")
	}

	if password == "" {
		return errorpkg.New("password is required")
	}

	if len(password) < 8 {
		return errorpkg.New("password length should be 8 characters at least")
	}

	a, err := mail.ParseAddress(email)
	if err != nil || a.Address != email {
		return errorpkg.New("email validation failed")
	}

	return nil
}
