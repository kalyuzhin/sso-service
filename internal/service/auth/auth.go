package auth

import (
	"context"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/kalyuzhin/sso-service/internal/config"
	errorpkg "github.com/kalyuzhin/sso-service/internal/error"
	"github.com/kalyuzhin/sso-service/internal/lib/jwt"
	"github.com/kalyuzhin/sso-service/internal/lib/refreshtoken"
	"github.com/kalyuzhin/sso-service/internal/model"
)

type Auth struct {
	userSaver    userSaver
	userProvider userProvider
	sessionSaver sessionSaver
	appProvider  appProvider
	cfg          config.Config
}

type sessionSaver interface {
	SaveRefreshSession(ctx context.Context, refreshTokenHash []byte, userID int64, ip, userAgent string,
		exparation time.Time) error
}

type userSaver interface {
	SaveUser(ctx context.Context, email string, passwordHash []byte) (userID int64, err error)
}

type userProvider interface {
	GetUser(ctx context.Context, email string) (user model.DBUser, err error)
}

type appProvider interface {
	App(ctx context.Context, appID int32) (a model.App, err error)
}

// New – ...
func New(provider userProvider, saver userSaver, provider2 appProvider, sessionSaver sessionSaver,
	cfg config.Config) *Auth {
	return &Auth{
		userSaver:    saver,
		userProvider: provider,
		appProvider:  provider2,
		sessionSaver: sessionSaver,
		cfg:          cfg,
	}
}

// Login – ...
func (a *Auth) Login(ctx context.Context, email, pswd string, appID int32, params model.UserRequestParams) (accessToken string, err error) {
	user, err := a.userProvider.GetUser(ctx, email)
	if err != nil {
		return accessToken, errorpkg.WrapErr(err, "can't get user from storage")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(pswd))
	if err != nil {
		return accessToken, errorpkg.WrapErr(err, "stored hash and password hash aren't equal")
	}

	app, err := a.appProvider.App(ctx, appID)
	if err != nil {
		return accessToken, errorpkg.WrapErr(err, "can't get app from storage")
	}

	accessToken, err = jwt.GenerateToken(app, user, time.Hour)
	if err != nil {
		return accessToken, errorpkg.WrapErr(err, "can't generate access token")
	}
	_, refreshTokenHash, err := refreshtoken.GenerateRefreshToken()
	if err != nil {
		return accessToken, errorpkg.WrapErr(err, "can't generate refresh token")
	}

	refreshTokenExpiration := time.Now().UTC().Add(a.cfg.RefreshTokenExparation)

	err = a.sessionSaver.SaveRefreshSession(ctx, refreshTokenHash, user.ID, params.IP, params.UserAgent, refreshTokenExpiration)
	if err != nil {
		return accessToken, errorpkg.WrapErr(err, "can't save refresh session")
	}

	return accessToken, nil
}

// Register – ...
func (a *Auth) Register(ctx context.Context, email, pswd string) (userID int64, err error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pswd), bcrypt.DefaultCost)
	if err != nil {
		return 0, errorpkg.WrapErr(err, "can't create hash from password")
	}

	userID, err = a.userSaver.SaveUser(ctx, email, hash)
	if err != nil {
		return 0, err
	}

	return userID, nil
}
