package auth

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	errorpkg "github.com/kalyuzhin/sso-service/internal/error"
	"github.com/kalyuzhin/sso-service/internal/model"
)

type Auth struct {
	userSaver    userSaver
	userProvider userProvider
	appProvider  appProvider
}

type userSaver interface {
	SaveUser(ctx context.Context, email string, passwordHash []byte) (userID int64, err error)
}

type userProvider interface {
	User(ctx context.Context, email string) (user model.DBUser, err error)
}

type appProvider interface {
	App(ctx context.Context, appID int32) (a model.App, err error)
}

// New – ...
func New(provider userProvider, saver userSaver, provider2 appProvider) *Auth {
	return &Auth{
		userSaver:    saver,
		userProvider: provider,
		appProvider:  provider2,
	}
}

// Login – ...
func (a *Auth) Login(ctx context.Context, email, pswd string, appID int32) (token string, err error) {
	user, err := a.userProvider.User(ctx, email)
	if err != nil {
		return token, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(pswd), bcrypt.DefaultCost)
	if err != nil {
		return token, errorpkg.WrapErr(err, "can't create hash from password")
	}

	err = bcrypt.CompareHashAndPassword(user.PasswordHash, hash)
	if err != nil {
		return token, errorpkg.WrapErr(err, "stored hash and password hash aren't equal")
	}

	return token, nil
}

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
