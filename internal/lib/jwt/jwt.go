package jwt

import (
	"crypto/rsa"
	"time"

	"github.com/golang-jwt/jwt/v5"

	errorpkg "github.com/kalyuzhin/sso-service/internal/error"
	"github.com/kalyuzhin/sso-service/internal/model"
)

const (
	kid = "ab1c1"
)

// GenerateToken – ...
func GenerateToken(app model.App, user model.DBUser, ttl time.Duration, privateKey *rsa.PrivateKey) (token string, err error) {
	tokenObj := jwt.New(jwt.SigningMethodRS256)
	// TODO: вынести генерацию kid и добавить его в конфиг
	tokenObj.Header["kid"] = kid

	claims := tokenObj.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(ttl).Unix()
	claims["app_id"] = app.ID

	token, err = tokenObj.SignedString(privateKey)
	if err != nil {
		return token, errorpkg.WrapErr(err, "can't convert token object into string")
	}

	return token, nil
}
